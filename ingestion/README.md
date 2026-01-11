# Porcool Ingestion Service

A Go-based data ingestion service that synchronizes data from MongoDB to MariaDB for the Porcool financial management application. The service operates as a RabbitMQ consumer, processing ingestion messages that trigger document synchronization.

## Overview

This service runs as a RabbitMQ consumer that listens for ingestion messages on the `porcool-ingestion-non-relational-database-to-relational-database` queue. When a message is received, it fetches the corresponding document from MongoDB's `succesfully_ingested_firestore_docs` collection and synchronizes the referenced documents to MariaDB.

```mermaid
flowchart TB
    subgraph Queue["Message Queue"]
        RabbitMQ[(RabbitMQ)]
    end

    subgraph Sources["Data Sources"]
        MongoDB[(MongoDB)]
    end

    subgraph Ingestion["Ingestion Service"]
        Consumer[RabbitMQ Consumer]
        Processor[Message Processor]
        Transform[Data Transformer]
    end

    subgraph Target["Target Database"]
        MariaDB[(MariaDB)]
    end

    RabbitMQ --> Consumer
    Consumer --> Processor
    Processor --> MongoDB
    MongoDB --> Transform
    Transform --> MariaDB
    MariaDB --> |Mark as synced| MongoDB
```

## Architecture

```mermaid
flowchart LR
    subgraph Container["Docker Container (scratch)"]
        Binary[Go Binary]
    end

    subgraph External["External Services (Production)"]
        RabbitMQ[(RabbitMQ)]
        MongoDB[(MongoDB)]
        MariaDB[(MariaDB)]
    end

    RabbitMQ --> Binary
    Binary <--> MongoDB
    Binary <--> MariaDB
```

## Features

- **RabbitMQ Consumer**: Event-driven message processing from the queue
- **Automatic Schema Migration**: Creates all required tables on startup if they don't exist
- **Domain Seeding**: Automatically populates domain/lookup tables with initial values
- **Idempotent Sync**: Uses upsert operations to safely handle duplicate syncs
- **Message-Based Processing**: Processes specific documents referenced in queue messages
- **Minimal Docker Image**: Uses multi-stage build with `scratch` base for tiny image size (~10MB)
- **Graceful Shutdown**: Properly handles SIGINT/SIGTERM signals and message acknowledgment
- **Centralized Logging**: Optional OpenSearch logging with 90-day retention and automatic fallback to stdout

## Message Flow

```mermaid
sequenceDiagram
    participant Queue as RabbitMQ Queue
    participant Consumer as Consumer
    participant MongoDB
    participant MariaDB

    Queue->>Consumer: Message: {successfullyIngestedFirestoreDocsID: "doc123"}
    Consumer->>MongoDB: Fetch succesfully_ingested_firestore_docs document
    MongoDB-->>Consumer: Document with map_collection_to_docs

    loop For each collection in map_collection_to_docs
        Consumer->>MongoDB: Fetch documents by IDs
        MongoDB-->>Consumer: Documents
        Consumer->>MariaDB: Upsert documents
        Consumer->>MongoDB: Mark documents as synced
    end

    Consumer->>Queue: Acknowledge message
```

## Queue Message Format

The service listens on the `porcool-ingestion-non-relational-database-to-relational-database` queue and expects messages in the following format:

```json
{
  "successfullyIngestedFirestoreDocsID": "string"
}
```

The `successfullyIngestedFirestoreDocsID` value is the `_id` of a document in the `succesfully_ingested_firestore_docs` MongoDB collection. This document contains a `map_collection_to_docs` field that maps collection names to arrays of document IDs to be synced.

### Example `succesfully_ingested_firestore_docs` Document

```json
{
  "_id": "doc123",
  "createdAt": "2024-01-15T10:30:00Z",
  "ingestedBy": "firestore-sync-service",
  "map_collection_to_docs": {
    "users": ["user1", "user2"],
    "expenses": ["expense1", "expense2", "expense3"],
    "banks": ["bank1"]
  },
  "onPremiseSyncService": "porcool-ingestion-non-relational-database-to-relational-database"
}
```

## Database Schema

```mermaid
erDiagram
    domain {
        bigint id PK
        varchar guid UK
        varchar name
        varchar type
        varchar source
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    user {
        bigint id PK
        varchar guid UK
        varchar source_id
        varchar first_name
        varchar last_name
        varchar email UK
        boolean fl_admin
        decimal monthly_income
        boolean fl_payment_requested
        boolean fl_payment_pending
        boolean fl_payment_paid
        varchar current_spending_date
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    financial_institution {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        varchar name
        boolean fl_credit_card
        boolean fl_money_movement
        boolean fl_investment
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    expense {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        varchar spending_date__YYYY_MM
        bigint id_status FK
        bigint id_type FK
        date validity_period_date
        boolean fl_indeterminate_validity_period_date
        varchar name
        decimal total_amount
        decimal total_paid_amount
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    expense_installment {
        bigint id PK
        varchar guid UK
        bigint expense_id FK
        decimal amount
        decimal paid_amount
        bigint id_status FK
        date due_date
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    expense_automatic_workflow {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        longtext base64_image
        text description
        longtext extracted_expense_content_from_image
        varchar spending_date__YYYY_MM
        timestamp sync_processed_date
        bigint id_sync_status FK
        text processing_message
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    expense_automatic_workflow_pre_saved_description {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        text description
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    additional_balance {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        varchar spending_date__YYYY_MM
        decimal amount
        text description
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    balance_history {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        varchar spending_date__YYYY_MM
        decimal amount
        decimal last_month_amount
        decimal monthly_income
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    service_payment {
        bigint id PK
        varchar guid UK
        varchar source_id
        bigint user_id FK
        date service_payment_date
        bigint service_payment_type_id FK
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    system_settings {
        bigint id PK
        varchar guid UK
        varchar source_id
        boolean fl_block_user_registration
        boolean fl_maintenance
        json json_sync_metadata
        timestamp created_at
        varchar created_by
        timestamp updated_at
        varchar updated_by
    }

    user ||--o{ financial_institution : "has"
    user ||--o{ expense : "has"
    user ||--o{ expense_automatic_workflow : "has"
    user ||--o{ expense_automatic_workflow_pre_saved_description : "has"
    user ||--o{ additional_balance : "has"
    user ||--o{ balance_history : "has"
    user ||--o{ service_payment : "has"
    expense ||--o{ expense_installment : "has"
    domain ||--o{ expense : "id_status"
    domain ||--o{ expense : "id_type"
    domain ||--o{ expense_installment : "id_status"
    domain ||--o{ expense_automatic_workflow : "id_sync_status"
    domain ||--o{ service_payment : "service_payment_type_id"
```

## Domain Values (Seeded on Startup)

| Source | Type | Values |
|--------|------|--------|
| expense_automatic_workflow | id_sync_status | pending, success, error |
| expense | id_status | pending, partially_paid, paid |
| expense | id_type | expense, invoice, savings |
| expense_installment | id_status | pending, partially_paid, paid |
| service_payment | service_payment_type_id | PayPal |

## MongoDB Collections

The service reads from the following MongoDB collections (originally synced from Firestore):

### Collection: `succesfully_ingested_firestore_docs`

This is the main collection that drives the ingestion process. Documents contain:

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `createdAt` | date | Creation timestamp |
| `ingestedBy` | string | Service that created the document |
| `map_collection_to_docs` | object | Map of collection names to document IDs |
| `onPremiseSyncService` | string | Service identifier for on-premise sync |

### Collection: `users`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID (same as Firestore doc ID) |
| `_firestoreCreateTime` | string | Firestore creation timestamp |
| `_firestorePath` | string | Original Firestore document path |
| `_firestoreUpdateTime` | string | Firestore update timestamp |
| `_importedAt` | date | MongoDB import timestamp |
| `admin` | bool | Admin flag |
| `email` | string | User email |
| `lastName` | string | User last name |
| `lookingAtSpendingDate` | string | Current spending date (YYYY-MM) |
| `monthlyIncome` | int | Monthly income |
| `name` | string | User first name |
| `onPremiseSyncDatetime` | null/date | Sync timestamp (null = pending) |
| `onPremiseSyncService` | null/string | Service that synced the document |
| `paidPayment` | bool | Payment paid flag |
| `requestedPayment` | bool | Payment requested flag |
| `pendingPayment` | bool | Payment pending flag |

### Collection: `expenses`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `alreadyPaidAmount` | int | Amount already paid |
| `amount` | int/double | Total expense amount |
| `created` | string | Creation timestamp |
| `expenseName` | string | Expense name |
| `indeterminateValidity` | bool | Indeterminate validity flag |
| `source` | string | Expense source |
| `spendingDate` | string | Current spending date |
| `status` | string | Expense status (pending/partially_paid/paid) |
| `type` | string | Expense type (expense/invoice/savings) |
| `updated` | string | Update timestamp |
| `user` | string | User ID reference |
| `validity` | null/string | Validity period date |

### Collection: `banks` (Financial Institutions)

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `cartaoCredito` | bool | Credit card flag |
| `created` | string | Creation timestamp |
| `investimentos` | bool | Investment flag |
| `movimentacaoDinheiro` | bool | Money movement flag |
| `nome` | string | Bank/institution name |
| `observacoes` | string | Notes/observations |
| `updated` | string | Update timestamp |
| `user` | string | User ID reference |

### Collection: `additional_balances`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `balance` | int/double | Balance amount |
| `created` | string | Creation timestamp |
| `description` | string | Balance description |
| `spendingDate` | string | Current spending date |
| `user` | string | User ID reference |

### Collection: `balance_history`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `balance` | int/double | Current balance |
| `created` | string | Creation timestamp |
| `lastMonthBalance` | int/double | Last month's balance |
| `monthlyIncome` | int | Monthly income |
| `spendingDate` | string | Current spending date |
| `user` | string | User ID reference |

### Collection: `expense_automatic_workflow`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `base64_image` | string | Base64 encoded image |
| `created` | string | Creation timestamp |
| `description` | string | Workflow description |
| `extracted_expense_content_from_image` | array/string | Extracted expense data |
| `processingMessage` | string | Processing status message |
| `spendingDate` | string | Current spending date |
| `syncProcessedDate` | string | Sync processed timestamp |
| `syncStatus` | string | Sync status (pending/success/error) |
| `user` | string | User ID reference |

### Collection: `expense_automatic_workflow_pre_saved_description`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `created` | string | Creation timestamp |
| `description` | string | Pre-saved description text |
| `user` | string | User ID reference |

### Collection: `payments`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `paymentDate` | string | Payment date |
| `user` | string | User ID reference |

### Collection: `settings`

| Field | Type | Description |
|-------|------|-------------|
| `_id` | string | Document ID |
| `blockUserRegistration` | bool | Block new registrations flag |
| `maintenance` | bool | Maintenance mode flag |
| `syncMetadata` | array | Sync metadata array |

## Field Mappings (MongoDB to MariaDB)

```mermaid
flowchart LR
    subgraph MongoDB["MongoDB Field Names"]
        m1["name"]
        m2["lastName"]
        m3["admin"]
        m4["lookingAtSpendingDate"]
        m5["requestedPayment"]
        m6["pendingPayment"]
        m7["paidPayment"]
        m8["expenseName"]
        m9["amount"]
        m10["alreadyPaidAmount"]
        m11["indeterminateValidity"]
        m12["nome"]
        m13["cartaoCredito"]
        m14["movimentacaoDinheiro"]
        m15["investimentos"]
        m16["balance"]
        m17["lastMonthBalance"]
    end

    subgraph MariaDB["MariaDB Field Names"]
        d1["first_name"]
        d2["last_name"]
        d3["fl_admin"]
        d4["current_spending_date"]
        d5["fl_payment_requested"]
        d6["fl_payment_pending"]
        d7["fl_payment_paid"]
        d8["name"]
        d9["total_amount"]
        d10["total_paid_amount"]
        d11["fl_indeterminate_validity_period_date"]
        d12["name"]
        d13["fl_credit_card"]
        d14["fl_money_movement"]
        d15["fl_investment"]
        d16["amount"]
        d17["last_month_amount"]
    end

    m1 --> d1
    m2 --> d2
    m3 --> d3
    m4 --> d4
    m5 --> d5
    m6 --> d6
    m7 --> d7
    m8 --> d8
    m9 --> d9
    m10 --> d10
    m11 --> d11
    m12 --> d12
    m13 --> d13
    m14 --> d14
    m15 --> d15
    m16 --> d16
    m17 --> d17
```

## Configuration

The service is configured via environment variables:

### Database Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MARIADB_HOST` | MariaDB hostname | `localhost` |
| `MARIADB_PORT` | MariaDB port | `3306` |
| `MARIADB_USER` | MariaDB username | `root` |
| `MARIADB_PASSWORD` | MariaDB password | `` |
| `MARIADB_DATABASE` | MariaDB database name | `porcool` |
| `MONGODB_URI` | MongoDB connection URI | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | MongoDB database name | `porcool` |

### Message Queue Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `RABBITMQ_URI` | RabbitMQ connection URI | `amqp://guest:guest@localhost:5672/` |
| `RABBITMQ_QUEUE_NAME` | RabbitMQ queue name | `porcool-ingestion-non-relational-database-to-relational-database` |
| `INGESTION_BATCH_SIZE` | Max documents per sync batch | `100` |

### OpenSearch Logging Configuration

The service supports centralized logging to OpenSearch with automatic fallback to stdout if OpenSearch is unavailable.

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENSEARCH_ENABLED` | Enable OpenSearch logging | `false` |
| `OPENSEARCH_URL` | OpenSearch URL (e.g., `https://opensearch:9200`) | `` |
| `OPENSEARCH_USERNAME` | OpenSearch username for authentication | `` |
| `OPENSEARCH_PASSWORD` | OpenSearch password for authentication | `` |
| `OPENSEARCH_INDEX_PREFIX` | Prefix for log indices | `porcool-ingestion-non-relational-database-to-relational-database` |
| `OPENSEARCH_RETENTION_DAYS` | Log retention period in days | `90` |

#### OpenSearch Logging Features

- **Automatic Fallback**: If OpenSearch is not configured or unavailable, logs are written to stdout
- **Automatic Reconnection**: The service periodically attempts to reconnect if OpenSearch becomes unavailable
- **Index State Management**: Automatically creates ISM policies for log retention (90 days by default)
- **Daily Indices**: Logs are stored in daily indices (e.g., `porcool-ingestion-non-relational-database-to-relational-database-2024.01.15`)
- **Structured Logging**: All logs include timestamp, level, message, service name, and optional fields

#### Log Entry Format

```json
{
  "@timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "message": "Processing ingestion message for document ID: doc123",
  "service": "porcool-ingestion-non-relational-database-to-relational-database",
  "host": "ingestion-pod-1",
  "fields": {
    "docID": "doc123",
    "collections": 5
  }
}
```

## Project Structure

```
ingestion/
├── main.go                              # Application entry point
├── go.mod                               # Go module definition
├── go.sum                               # Dependency checksums
├── Dockerfile                           # Multi-stage Docker build
├── docker-compose.yml                   # Development environment
├── .dockerignore                        # Docker build exclusions
├── README.md                            # This file
└── internal/
    ├── config/
    │   ├── config.go                    # Configuration loading
    │   └── config_test.go               # Config tests
    ├── models/
    │   ├── models.go                    # Data models and domain seeds
    │   └── models_test.go               # Model tests
    ├── database/
    │   ├── mariadb/
    │   │   ├── connection.go            # MariaDB connection and migrations
    │   │   ├── connection_test.go       # Connection tests
    │   │   ├── repository.go            # Database repositories
    │   │   └── repository_test.go       # Repository tests
    │   └── mongodb/
    │       ├── connection.go            # MongoDB connection and queries
    │       └── connection_test.go       # MongoDB tests
    ├── logging/
    │   ├── logger.go                    # Logger with fallback support
    │   ├── logger_test.go               # Logger tests
    │   ├── opensearch.go                # OpenSearch client
    │   └── opensearch_test.go           # OpenSearch tests
    ├── queue/
    │   └── rabbitmq/
    │       ├── consumer.go              # RabbitMQ consumer
    │       └── consumer_test.go         # Consumer tests
    └── ingestion/
        ├── service.go                   # Main ingestion service
        └── service_test.go              # Service tests
```

## Running Locally

### Prerequisites

- Go 1.22+
- Docker and Docker Compose (for running databases and RabbitMQ)
- Make (optional)

### Development Setup

**Note:** The docker-compose.yml file is for local testing only. In production, all services (MariaDB, MongoDB, RabbitMQ) run externally.

1. Start the test infrastructure:
```bash
docker-compose up -d mariadb mongodb rabbitmq
```

2. Run the service:
```bash
go run main.go
```

### Running Tests

```bash
go test ./...
```

### Running with Coverage

```bash
go test -cover ./...
```

## Docker

### Building the Image

```bash
docker build -t porcool-ingestion-non-relational-database-to-relational-database .
```

### Running with Docker Compose

```bash
docker-compose up -d
```

### Image Size

The final Docker image is built using a multi-stage build with a `scratch` base, resulting in a minimal image size of approximately **10-15MB**.

```mermaid
flowchart LR
    subgraph Build["Build Stage (~1GB)"]
        Go[golang:1.22-alpine]
        Deps[Dependencies]
        Compile[Compile Binary]
    end

    subgraph Final["Final Stage (~10MB)"]
        Scratch[scratch]
        Binary[Go Binary]
        Certs[CA Certificates]
    end

    Go --> Deps --> Compile --> Binary
    Scratch --> Binary
    Scratch --> Certs
```

## Startup Process

```mermaid
flowchart TD
    Start([Start]) --> LoadConfig[Load Configuration]
    LoadConfig --> InitLogger[Initialize Logger]
    InitLogger --> |OpenSearch enabled| ConnectOpenSearch[Connect to OpenSearch]
    InitLogger --> |OpenSearch disabled| UseStdout[Use stdout logging]
    ConnectOpenSearch --> |Success or Failure| UseStdout
    UseStdout --> ConnectMaria[Connect to MariaDB]
    ConnectMaria --> |Success| RunMigrations[Run Migrations]
    ConnectMaria --> |Failure| Exit1([Exit with Error])
    RunMigrations --> |Success| SeedDomains[Seed Domains]
    RunMigrations --> |Failure| Exit2([Exit with Error])
    SeedDomains --> |Success| ConnectMongo[Connect to MongoDB]
    SeedDomains --> |Failure| Exit3([Exit with Error])
    ConnectMongo --> |Success| ConnectRabbitMQ[Connect to RabbitMQ]
    ConnectMongo --> |Failure| Exit4([Exit with Error])
    ConnectRabbitMQ --> |Success| StartConsumer[Start RabbitMQ Consumer]
    ConnectRabbitMQ --> |Failure| Exit5([Exit with Error])
    StartConsumer --> WaitSignal[Wait for Shutdown Signal]
    WaitSignal --> |SIGINT/SIGTERM| Shutdown[Graceful Shutdown]
    Shutdown --> End([End])
```

## Message Processing Flow

```mermaid
flowchart TD
    Receive[Receive Message from Queue] --> Parse[Parse JSON Message]
    Parse --> |Invalid JSON| Reject[Reject Message]
    Parse --> |Valid JSON| Validate[Validate Message]
    Validate --> |Missing ID| Reject
    Validate --> |Valid| Fetch[Fetch succesfully_ingested_firestore_docs]
    Fetch --> |Not Found| Ack[Acknowledge Message]
    Fetch --> |Found| ProcessCollections[Process map_collection_to_docs]
    ProcessCollections --> |For each collection| SyncDocs[Sync Documents to MariaDB]
    SyncDocs --> MarkSynced[Mark Documents as Synced in MongoDB]
    MarkSynced --> Ack
```

## Supported Collections

The service currently supports ingestion for the following collections:

| Collection Name | Target Table |
|----------------|--------------|
| `users` | `user` |
| `expenses` | `expense`, `expense_installment` |
| `banks` | `financial_institution` |
| `additional_balances` | `additional_balance` |
| `balance_history` | `balance_history` |
| `expense_automatic_workflow` | `expense_automatic_workflow` |
| `expense_automatic_workflow_pre_saved_description` | `expense_automatic_workflow_pre_saved_description` |
| `payments` | `service_payment` |
| `settings` | `system_settings` |

## Collection Ingestion Order

Collections are processed in a specific order to ensure referential integrity:

1. **users** - Must be processed first since almost all other collections depend on user references
2. **banks** - Financial institutions
3. **expenses** - Expense records (with special handling for invoice/savings)
4. **additional_balances** - Additional balance entries
5. **balance_history** - Balance history records
6. **expense_automatic_workflow** - Automatic expense workflows
7. **expense_automatic_workflow_pre_saved_description** - Pre-saved descriptions
8. **payments** - Service payments
9. **settings** - System settings

## Expense Ingestion Rules

The expense collection has special ingestion rules based on the expense type:

### Simple Expense (type = "expense")

When the expense type is `"expense"`, the service creates a single record in the `expense` table with all fields populated.

### Invoice/Savings Without Validity

When the expense type is `"invoice"` or `"savings"` and the `validity` field is NULL or empty, the service creates a single record in the `expense` table (same as simple expense).

### Invoice/Savings With Validity (Aggregation)

When the expense type is `"invoice"` or `"savings"` and the `validity` field has a value:

1. **Aggregate Expense Record**: Creates a single "generic" expense record in the `expense` table with:
   - Empty `spending_date__YYYY_MM` (indicates aggregate)
   - NULL `id_status` (not applicable for aggregate)
   - `validity_period_date` set to the validity date
   - `name` from the expense
   - `total_amount` = 0 (not tracked for aggregates)
   - `total_paid_amount` = 0 (not tracked for aggregates)

2. **Installment Generation**: Creates records in `expense_installment` table:
   - One installment for each MongoDB expense record with the same `expenseName`, `user`, and `validity`
   - Additional "pending" installments generated for each month from the last existing expense date + 1 month until the validity date

### Duplicate Prevention (IMPORTANT)

**Critical for invoice/savings expenses with validity:**

- The service checks for existing aggregate expense records by matching `name`, `validity_period_date`, and `user_id` with empty `spending_date`
- Before creating an installment, it checks if one already exists for that expense and due date
- This prevents duplicate records when the same expense is processed multiple times
- Both checks are performed using dedicated repository methods:
  - `GetExpenseByNameValidityUser()` for expense aggregates
  - `GetInstallmentByExpenseAndDate()` for installments

```mermaid
flowchart TD
    Expense[Expense Document] --> CheckType{Expense Type?}

    CheckType --> |expense| SimpleExpense[Create Simple Expense Record]
    CheckType --> |invoice/savings| CheckValidity{Has Validity?}

    CheckValidity --> |No validity| SimpleExpense
    CheckValidity --> |Has validity| CheckExisting{Existing Aggregate?}

    CheckExisting --> |Yes| UseExisting[Use Existing Expense ID]
    CheckExisting --> |No| CreateAggregate[Create Aggregate Expense]

    UseExisting --> FetchRelated[Fetch All Related Expenses]
    CreateAggregate --> FetchRelated

    FetchRelated --> CreateInstallments[Create/Update Installments for Each]
    CreateInstallments --> GeneratePending[Generate Pending Installments Until Validity]

    SimpleExpense --> MarkSynced[Mark as Synced in MongoDB]
    GeneratePending --> MarkSynced
```

## Spending Date Format

All spending dates are converted to `YYYY/MM` format before being inserted into MariaDB. The service handles the following input formats:

| Input Format | Example | Output |
|--------------|---------|--------|
| YYYYMM | 202312 | 2023/12 |
| YYYY-MM | 2023-12 | 2023/12 |
| YYYY/MM | 2023/12 | 2023/12 |

## MongoDB Sync Fields

When a document is successfully synced to MariaDB, the service marks it with the following fields in MongoDB:

| Field | Type | Description |
|-------|------|-------------|
| `onPremiseRelationalDBSyncDatetime` | Date | Timestamp when the document was synced |
| `onPremiseRelationalDBSyncService` | String | Service name: `porcool-ingestion-non-relational-database-to-relational-database` |

## Audit Fields

All records in MariaDB include audit fields:

| Field | Description |
|-------|-------------|
| `created_at` | Timestamp when the record was created (current datetime) |
| `created_by` | Service name: `porcool-ingestion-non-relational-database-to-relational-database` |
| `updated_at` | Timestamp when the record was last updated (current datetime on update) |
| `updated_by` | Service name: `porcool-ingestion-non-relational-database-to-relational-database` |

## License

This project is part of the Porcool application.
