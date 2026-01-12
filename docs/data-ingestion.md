# Data Ingestion Pipeline

This document describes the two-stage data ingestion pipeline that synchronizes data from **Firebase Firestore** to **MariaDB**.

## Overview

The ingestion pipeline consists of two independent services that work together through a message queue:

```mermaid
flowchart LR
    subgraph Databases
        FS[(Firestore)]
        MG[(MongoDB)]
        MDB[(MariaDB)]
    end

    subgraph Infrastructure
        MQ{{RabbitMQ}}
    end

    subgraph Services
        N8N[N8N Flow<br/>Stage 1]
        GO[Go Ingestion<br/>Stage 2]
    end

    FS -->|Fetch unsynced docs| N8N
    N8N -->|Insert documents| MG
    N8N -->|Publish message| MQ
    MQ -->|Consume message| GO
    MG -->|Fetch documents| GO
    GO -->|Transform & Upsert| MDB
```

## Stage 1: N8N Flow (Firestore to MongoDB)

**Location**: External n8n workflow
**Schedule**: Runs every 6 hours
**Purpose**: Fetches unsynced documents from Firestore and copies them to MongoDB

### Workflow Nodes

```mermaid
flowchart LR
    A[Every 6 Hours] --> B[Config]
    B --> C[Fetch Docs From<br/>All Collections]
    C --> D[Replicate Docs<br/>To MongoDB]
    D --> E[Track Ingested<br/>Docs in MongoDB]
    E --> F[Update Firestore Docs<br/>With Metadata]
    F --> G[Update Settings<br/>Sync Metadata]
    G --> H[Publish RabbitMQ<br/>Message]
```

| Node | Description |
|------|-------------|
| Every 6 Hours | Schedule trigger that initiates the flow |
| Config | Configuration node with Firebase SA, MongoDB URI, and sync settings |
| Fetch Docs From All Collections | Queries Firestore for documents where `onPremiseSyncDatetime IS NULL` |
| Replicate Docs To MongoDB | Inserts fetched documents into MongoDB collections |
| Track Ingested Docs in MongoDB | Creates a tracking record with all successfully synced document IDs |
| Update Firestore Docs With "onPremise*" Metadata | Marks Firestore documents as synced |
| Update Settings Sync Metadata | Updates the `syncMetadata` array in settings collection |
| Publish RabbitMQ Message | Notifies the Go service that new documents are ready |

### Collections Synced

- `additional_balances`
- `balance_history`
- `banks`
- `expense_automatic_workflow`
- `expense_automatic_workflow_pre_saved_description`
- `expenses`
- `payments`
- `settings`
- `users`

### Firestore Query Logic

Documents are fetched using the following criteria:

```javascript
{
  where: {
    unaryFilter: {
      field: { fieldPath: 'onPremiseSyncDatetime' },
      op: 'IS_NULL'
    }
  },
  limit: 150  // Configurable via DOCS_LIMIT_PER_COLLECTION
}
```

Collections with a `created` field are ordered by `created DESC` (excludes: `users`, `payments`, `settings`).

### MongoDB Document Structure

Each document inserted into MongoDB includes:

| Field | Description |
|-------|-------------|
| `_id` | Original Firestore document ID |
| `...data` | All original Firestore document fields |
| `_firestorePath` | Full path to the Firestore document |
| `_firestoreCreateTime` | Original creation timestamp |
| `_firestoreUpdateTime` | Original update timestamp |
| `_importedAt` | Timestamp when imported to MongoDB |

### Tracking Document

After successful MongoDB insertion, a tracking document is created in `succesfully_ingested_firestore_docs`:

```javascript
{
  onPremiseSyncService: 'porcool-raw-firestore-copy-ingestion',
  map_collection_to_docs: {
    'users': ['doc_id_1', 'doc_id_2'],
    'expenses': ['doc_id_3', 'doc_id_4'],
    // ... other collections
  },
  ingestedBy: null,  // Set by Go service after processing
  createdAt: ISODate()
}
```

### RabbitMQ Message

The message published to RabbitMQ:

```json
{
  "successfullyIngestedFirestoreDocsID": "<tracking_document_id>"
}
```

**Queue**: `porcool-ingestion-non-relational-database-to-relational-database`

### Firestore Metadata Updates

After successful sync, Firestore documents are updated with:

```javascript
{
  onPremiseSyncDatetime: <current_timestamp>,
  onPremiseSyncService: 'porcool-raw-firestore-copy-ingestion'
}
```

---

## Stage 2: Go Ingestion Service (MongoDB to MariaDB)

**Location**: `/ingestion/` directory
**Trigger**: RabbitMQ message consumer
**Purpose**: Transforms and loads MongoDB documents into normalized MariaDB tables

### Service Architecture

```mermaid
flowchart TB
    subgraph Go Ingestion Service
        RMQ[RabbitMQ Consumer]
        IS[Ingestion Service]
        RL[Repository Layer]
    end

    MG[(MongoDB)]
    MDB[(MariaDB)]

    RMQ -->|Message received| IS
    IS -->|Query documents| MG
    MG -->|Document data| IS
    IS -->|Transform| RL
    RL -->|Upsert| MDB
```

### Processing Flow

```mermaid
sequenceDiagram
    participant RMQ as RabbitMQ
    participant GO as Go Service
    participant MG as MongoDB
    participant MDB as MariaDB

    RMQ->>GO: Message with tracking ID
    GO->>MG: Fetch tracking document
    MG-->>GO: Tracking document

    loop For each collection
        GO->>MG: Fetch documents by IDs
        MG-->>GO: Document data
        GO->>GO: Transform data
        GO->>MDB: Upsert records
    end

    GO->>MG: Update tracking document
    GO->>MG: Update sync metadata
    GO->>RMQ: ACK message
```

1. **Message Received**: Consumer receives RabbitMQ message with tracking document ID
2. **Fetch Tracking Document**: Retrieves the tracking document from `succesfully_ingested_firestore_docs`
3. **Process Collections**: Iterates through each collection in dependency order
4. **Transform & Upsert**: For each document, transforms data and upserts to MariaDB
5. **Mark Processed**: Updates tracking document with `ingestedBy` and `ingestedAt`
6. **Update Sync Metadata**: Optionally updates Firestore settings with sync timestamp

### Collection Processing Order

Collections are processed in strict dependency order:

```mermaid
flowchart TD
    U[users] --> B[banks]
    U --> E[expenses]
    U --> AB[additional_balances]
    U --> BH[balance_history]
    U --> EAW[expense_automatic_workflow]
    U --> EAWPSD[expense_automatic_workflow_pre_saved_description]
    U --> P[payments]
    S[settings]
```

| Order | Collection | Dependency |
|-------|------------|------------|
| 1 | users | None |
| 2 | banks | users |
| 3 | expenses | users |
| 4 | additional_balances | users |
| 5 | balance_history | users |
| 6 | expense_automatic_workflow | users |
| 7 | expense_automatic_workflow_pre_saved_description | users |
| 8 | payments | users |
| 9 | settings | None |

### Data Transformations

#### Users

| MongoDB Field | MariaDB Field |
|---------------|---------------|
| `id` | `source_id` |
| `name` | `first_name` |
| `lastName` | `last_name` |
| `email` | `email` |
| `monthlyIncome` | `monthly_income` |
| `lookingAtSpendingDate` | `current_spending_date` (formatted) |
| `admin` | `fl_admin` |

#### Expenses

The expense transformation handles two distinct types:

```mermaid
flowchart TD
    E[Expense Document] --> T{Type?}
    T -->|expense| SE[Simple Expense]
    T -->|invoice/savings| AE{Has validity?}

    SE --> SR[Single record with<br/>spending_date populated]

    AE -->|Yes| AR[Parent record + Installments]
    AE -->|No / Indeterminate| IND[Single record<br/>No installments]

    AR --> I1[Installment 1]
    AR --> I2[Installment 2]
    AR --> IN[Installment N]
```

**Simple Expenses** (type: `expense`):
- Single record with `spending_date__YYYY_MM` populated
- Direct field mapping

**Aggregate Expenses** (type: `invoice` or `savings` with validity):
- Multiple MongoDB documents aggregated into one parent expense
- Parent record has empty `spending_date__YYYY_MM`
- Individual installments created in `expense_installment` table

**Indeterminate Validity**: Invoice/savings expenses may have indeterminate validity (no end date). In this case:
- Treated as a simple expense (single record)
- No installments are created
- The `fl_indeterminate_validity_period_date` flag is set to true

| MongoDB Field | MariaDB Field |
|---------------|---------------|
| `id` | `source_id` |
| `name` | `name` |
| `amount` | `total_amount` |
| `alreadyPaidAmount` | `total_paid_amount` |
| `status` | `id_status` (domain lookup) |
| `type` | `id_type` (domain lookup) |
| `validity` | `validity_period_date` |
| `indeterminateValidity` | `fl_indeterminate_validity_period_date` |
| `spendingDate` | `spending_date__YYYY_MM` |

**Installment Generation Logic** (only for invoice/savings with validity):
1. Find all MongoDB expenses with same name, validity, and user
2. Create installment for each existing expense record
3. Generate pending installments for remaining months until validity date

#### Financial Institutions (Banks)

| MongoDB Field | MariaDB Field |
|---------------|---------------|
| `id` | `source_id` |
| `nome` | `name` |
| `cartaoCredito` | `fl_credit_card` |
| `movimentacaoDinheiro` | `fl_money_movement` |
| `investimentos` | `fl_investment` |
| `user` | `user_id` (lookup) |

### Domain Values

Pre-seeded lookup values in the `domain` table:

| Source | Type | Values |
|--------|------|--------|
| expense_automatic_workflow | id_sync_status | pending, success, error |
| expense | id_status | pending, partially_paid, paid |
| expense | id_type | expense, invoice, savings |
| expense_installment | id_status | pending, partially_paid, paid |
| service_payment | service_payment_type_id | PayPal |

### Upsert Strategy

All records use an upsert pattern:

1. Look up existing record by `source_id` (MongoDB document ID)
2. If not found: INSERT with new UUID as `guid`
3. If found: UPDATE all fields
4. Set audit fields: `created_by`/`updated_by`, timestamps

### MongoDB Sync Tracking

After processing each document, MongoDB records are updated:

```javascript
{
  onPremiseRelationalDBSyncDatetime: <timestamp>,
  onPremiseRelationalDBSyncService: 'porcool-relational-db-ingestion'
}
```

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `MARIADB_HOST` | localhost | MariaDB host |
| `MARIADB_PORT` | 3306 | MariaDB port |
| `MARIADB_USER` | root | Database user |
| `MARIADB_PASSWORD` | (empty) | Database password |
| `MARIADB_DATABASE` | porcool | Database name |
| `MONGODB_URI` | mongodb://localhost:27017 | MongoDB connection |
| `MONGODB_DATABASE` | porcool | MongoDB database |
| `RABBITMQ_URI` | amqp://guest:guest@localhost:5672/ | RabbitMQ connection |
| `RABBITMQ_QUEUE_NAME` | porcool-ingestion-non-relational-database-to-relational-database | Queue name |
| `INGESTION_BATCH_SIZE` | 100 | Documents per query |
| `FIREBASE_ENABLED` | false | Enable Firestore sync metadata updates |

---

## Complete Data Flow

```mermaid
flowchart TB
    subgraph Stage1["Stage 1: N8N Flow (Every 6 hours)"]
        direction TB
        FS1[(Firestore)] -->|Query where<br/>syncDatetime IS NULL| N8N[N8N Workflow]
        N8N -->|Insert with metadata| MG1[(MongoDB)]
        N8N -->|Update sync fields| FS1
        MG1 -->|Create tracking doc| TD[Tracking Document]
        TD -->|Publish message| RMQ{{RabbitMQ}}
    end

    subgraph Stage2["Stage 2: Go Service (Event-driven)"]
        direction TB
        RMQ -->|Consume message| GO[Go Ingestion Service]
        GO -->|Fetch tracking doc| MG2[(MongoDB)]
        MG2 -->|Fetch documents by ID| GO
        GO -->|Transform & Upsert| MDB[(MariaDB)]
        GO -->|Update sync metadata| MG2
    end

    style Stage1 fill:#e1f5fe
    style Stage2 fill:#f3e5f5
```

### Detailed Flow Sequence

```mermaid
sequenceDiagram
    box Stage 1 - N8N Flow
        participant FS as Firestore
        participant N8N as N8N Workflow
        participant MG as MongoDB
    end

    box Message Queue
        participant RMQ as RabbitMQ
    end

    box Stage 2 - Go Service
        participant GO as Go Service
        participant MDB as MariaDB
    end

    Note over FS,N8N: Every 6 hours
    N8N->>FS: Query docs where onPremiseSyncDatetime IS NULL
    FS-->>N8N: Return up to 150 docs per collection

    N8N->>MG: Insert documents with metadata
    MG-->>N8N: Insertion result

    N8N->>MG: Create tracking document
    MG-->>N8N: Tracking document ID

    N8N->>FS: Update docs with onPremise* metadata
    N8N->>FS: Update settings syncMetadata

    N8N->>RMQ: Publish message with tracking ID

    Note over RMQ,GO: Event-driven
    RMQ->>GO: Deliver message

    GO->>MG: Fetch tracking document
    MG-->>GO: Return tracking document

    loop For each collection in order
        GO->>MG: Fetch documents by IDs
        MG-->>GO: Return documents
        GO->>GO: Transform data
        GO->>MDB: Upsert records
    end

    GO->>MG: Mark tracking doc as processed
    GO->>MG: Update document sync metadata
    GO->>RMQ: ACK message
```

---

## Error Handling

### N8N Flow
- Duplicate key errors handled gracefully (partial inserts recorded)
- Failed collections logged but don't stop other collections
- RabbitMQ message only published after successful sync

### Go Service

```mermaid
flowchart TD
    MSG[Message Received] --> PARSE{Valid JSON?}
    PARSE -->|No| NACK1[NACK without requeue<br/>Dead letter]
    PARSE -->|Yes| PROC[Process Collections]
    PROC --> ERR{Processing Error?}
    ERR -->|Yes| NACK2[NACK with requeue<br/>Retry]
    ERR -->|No| ACK[ACK<br/>Mark as processed]
```

- Invalid JSON messages: NACK without requeue (dead letter)
- Processing errors: NACK with requeue (retry)
- Collection errors: Logged but other collections continue
- Success: ACK and mark tracking document as processed

---

## Idempotency

The pipeline is designed to be idempotent:

```mermaid
flowchart LR
    subgraph Stage1[Stage 1 Idempotency]
        FS1[Firestore Doc] -->|Has onPremiseSyncDatetime?| CHECK1{Check}
        CHECK1 -->|Yes| SKIP1[Skip - Already synced]
        CHECK1 -->|No| SYNC1[Sync to MongoDB]
    end

    subgraph Stage2[Stage 2 Idempotency]
        MG1[MongoDB Doc] -->|Upsert by source_id| MDB1[(MariaDB)]
    end

    subgraph Tracking[Tracking Idempotency]
        TD1[Tracking Doc] -->|Has ingestedBy?| CHECK2{Check}
        CHECK2 -->|Yes| SKIP2[Skip - Already processed]
        CHECK2 -->|No| PROC2[Process]
    end
```

1. **Firestore → MongoDB**: Documents already synced have `onPremiseSyncDatetime` set and won't be re-fetched
2. **MongoDB → MariaDB**: Upsert operations use `source_id` for deduplication
3. **Tracking Documents**: `ingestedBy` field prevents reprocessing

---

## Monitoring

### Key Metrics to Watch
- Documents synced per collection per run
- RabbitMQ queue depth
- Processing time per message
- Error rates per collection

### Sync Status Fields

| Database | Field | Purpose |
|----------|-------|---------|
| Firestore | `onPremiseSyncDatetime` | Stage 1 completion marker |
| Firestore | `onPremiseSyncService` | Stage 1 service identifier |
| MongoDB | `onPremiseRelationalDBSyncDatetime` | Stage 2 completion marker |
| MongoDB | `onPremiseRelationalDBSyncService` | Stage 2 service identifier |
| MongoDB | `succesfully_ingested_firestore_docs.ingestedBy` | Tracking document processing status |
