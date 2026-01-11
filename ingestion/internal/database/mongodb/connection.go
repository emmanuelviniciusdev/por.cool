package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/porcool/ingestion/internal/config"
)

// Connection represents a MongoDB connection
type Connection struct {
	client   *mongo.Client
	database *mongo.Database
	cfg      config.MongoDBConfig
}

// NewConnection creates a new MongoDB connection
func NewConnection(cfg config.MongoDBConfig) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.URI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Connection{
		client:   client,
		database: client.Database(cfg.Database),
		cfg:      cfg,
	}, nil
}

// Close closes the MongoDB connection
func (c *Connection) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}

// Database returns the MongoDB database
func (c *Connection) Database() *mongo.Database {
	return c.database
}

// Collection returns a MongoDB collection
func (c *Connection) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// MongoDocument represents a generic MongoDB document
type MongoDocument struct {
	ID   string                 `bson:"_id,omitempty"`
	Data map[string]interface{} `bson:",inline"`
}

// UserDocument represents a user document from MongoDB (collection: users)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// admin, email, lastName, lookingAtSpendingDate, monthlyIncome, name, onPremiseSyncDatetime,
// onPremiseSyncService, paidPayment, requestedPayment, pendingPayment
type UserDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	Admin                 bool       `bson:"admin"`
	Email                 string     `bson:"email"`
	LastName              string     `bson:"lastName"`
	LookingAtSpendingDate string     `bson:"lookingAtSpendingDate"`
	MonthlyIncome         float64    `bson:"monthlyIncome"`
	Name                  string     `bson:"name"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	PaidPayment           bool       `bson:"paidPayment"`
	RequestedPayment      bool       `bson:"requestedPayment"`
	PendingPayment        bool       `bson:"pendingPayment"`
}

// ExpenseDocument represents an expense document from MongoDB (collection: expenses)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// alreadyPaidAmount, amount, created, expenseName, indeterminateValidity, onPremiseSyncDatetime,
// onPremiseSyncService, source, spendingDate, status, type, updated, user, validity
type ExpenseDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	AlreadyPaidAmount     float64    `bson:"alreadyPaidAmount"`
	Amount                float64    `bson:"amount"`
	Created               string     `bson:"created"`
	ExpenseName           string     `bson:"expenseName"`
	IndeterminateValidity bool       `bson:"indeterminateValidity"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	Source                string     `bson:"source"`
	SpendingDate          string     `bson:"spendingDate"`
	Status                string     `bson:"status"`
	Type                  string     `bson:"type"`
	Updated               string     `bson:"updated"`
	User                  string     `bson:"user"`
	Validity              *string    `bson:"validity"`
}

// FinancialInstitutionDocument represents a financial institution from MongoDB (collection: banks)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// cartaoCredito, created, investimentos, movimentacaoDinheiro, nome, observacoes,
// onPremiseSyncDatetime, onPremiseSyncService, updated, user
type FinancialInstitutionDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	CartaoCredito         bool       `bson:"cartaoCredito"`
	Created               string     `bson:"created"`
	Investimentos         bool       `bson:"investimentos"`
	MovimentacaoDinheiro  bool       `bson:"movimentacaoDinheiro"`
	Nome                  string     `bson:"nome"`
	Observacoes           string     `bson:"observacoes"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	Updated               string     `bson:"updated"`
	User                  string     `bson:"user"`
}

// AdditionalBalanceDocument represents an additional balance from MongoDB (collection: additional_balances)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// balance, created, description, onPremiseSyncDatetime, onPremiseSyncService, spendingDate, user
type AdditionalBalanceDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	Balance               float64    `bson:"balance"`
	Created               string     `bson:"created"`
	Description           string     `bson:"description"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	SpendingDate          string     `bson:"spendingDate"`
	User                  string     `bson:"user"`
}

// BalanceHistoryDocument represents a balance history record from MongoDB (collection: balance_history)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// balance, created, lastMonthBalance, monthlyIncome, onPremiseSyncDatetime, onPremiseSyncService,
// spendingDate, user
type BalanceHistoryDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	Balance               float64    `bson:"balance"`
	Created               string     `bson:"created"`
	LastMonthBalance      float64    `bson:"lastMonthBalance"`
	MonthlyIncome         float64    `bson:"monthlyIncome"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	SpendingDate          string     `bson:"spendingDate"`
	User                  string     `bson:"user"`
}

// ExpenseAutomaticWorkflowDocument represents an expense automatic workflow from MongoDB
// (collection: expense_automatic_workflow)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// base64_image, created, description, extracted_expense_content_from_image, onPremiseSyncDatetime,
// onPremiseSyncService, processingMessage, spendingDate, syncProcessedDate, syncStatus, user
type ExpenseAutomaticWorkflowDocument struct {
	ID                               string      `bson:"_id"`
	FirestoreCreateTime              string      `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath                    string      `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime              string      `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt                       time.Time   `bson:"_importedAt,omitempty"`
	Base64Image                      string      `bson:"base64_image"`
	Created                          string      `bson:"created"`
	Description                      string      `bson:"description"`
	ExtractedExpenseContentFromImage interface{} `bson:"extracted_expense_content_from_image"`
	OnPremiseSyncDatetime            *time.Time  `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService             *string     `bson:"onPremiseSyncService"`
	ProcessingMessage                string      `bson:"processingMessage"`
	SpendingDate                     string      `bson:"spendingDate"`
	SyncProcessedDate                string      `bson:"syncProcessedDate"`
	SyncStatus                       string      `bson:"syncStatus"`
	User                             string      `bson:"user"`
}

// ExpenseAutomaticWorkflowPreSavedDescriptionDocument represents pre-saved descriptions
// (collection: expense_automatic_workflow_pre_saved_description)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// created, description, onPremiseSyncDatetime, onPremiseSyncService, user
type ExpenseAutomaticWorkflowPreSavedDescriptionDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	Created               string     `bson:"created"`
	Description           string     `bson:"description"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	User                  string     `bson:"user"`
}

// ServicePaymentDocument represents a service payment from MongoDB (collection: payments)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// onPremiseSyncDatetime, onPremiseSyncService, paymentDate, user
type ServicePaymentDocument struct {
	ID                    string     `bson:"_id"`
	FirestoreCreateTime   string     `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string     `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string     `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time  `bson:"_importedAt,omitempty"`
	OnPremiseSyncDatetime *time.Time `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string    `bson:"onPremiseSyncService"`
	PaymentDate           string     `bson:"paymentDate"`
	User                  string     `bson:"user"`
}

// SettingsDocument represents system settings from MongoDB (collection: settings)
// MongoDB fields: _id, _firestoreCreateTime, _firestorePath, _firestoreUpdateTime, _importedAt,
// blockUserRegistration, maintenance, onPremiseSyncDatetime, onPremiseSyncService, syncMetadata
type SettingsDocument struct {
	ID                    string        `bson:"_id"`
	FirestoreCreateTime   string        `bson:"_firestoreCreateTime,omitempty"`
	FirestorePath         string        `bson:"_firestorePath,omitempty"`
	FirestoreUpdateTime   string        `bson:"_firestoreUpdateTime,omitempty"`
	ImportedAt            time.Time     `bson:"_importedAt,omitempty"`
	BlockUserRegistration bool          `bson:"blockUserRegistration"`
	Maintenance           bool          `bson:"maintenance"`
	OnPremiseSyncDatetime *time.Time    `bson:"onPremiseSyncDatetime"`
	OnPremiseSyncService  *string       `bson:"onPremiseSyncService"`
	SyncMetadata          []interface{} `bson:"syncMetadata"`
}

// SuccessfullyIngestedFirestoreDocsDocument represents ingestion tracking
// (collection: succesfully_ingested_firestore_docs)
// MongoDB fields: _id, createdAt, ingestedBy, map_collection_to_docs, onPremiseSyncService
type SuccessfullyIngestedFirestoreDocsDocument struct {
	ID                   interface{}            `bson:"_id"`
	CreatedAt            time.Time              `bson:"createdAt"`
	IngestedBy           *string                `bson:"ingestedBy"`
	MapCollectionToDocs  map[string]interface{} `bson:"map_collection_to_docs"`
	OnPremiseSyncService string                 `bson:"onPremiseSyncService"`
}

// GetPendingSyncUsers fetches users that need to be synced
func (c *Connection) GetPendingSyncUsers(ctx context.Context, limit int) ([]UserDocument, error) {
	collection := c.Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []UserDocument
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// GetPendingSyncExpenses fetches expenses that need to be synced
func (c *Connection) GetPendingSyncExpenses(ctx context.Context, limit int) ([]ExpenseDocument, error) {
	collection := c.Collection("expenses")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []ExpenseDocument
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, fmt.Errorf("failed to decode expenses: %w", err)
	}

	return expenses, nil
}

// GetPendingSyncFinancialInstitutions fetches financial institutions that need to be synced
func (c *Connection) GetPendingSyncFinancialInstitutions(ctx context.Context, limit int) ([]FinancialInstitutionDocument, error) {
	collection := c.Collection("banks")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find financial institutions: %w", err)
	}
	defer cursor.Close(ctx)

	var institutions []FinancialInstitutionDocument
	if err := cursor.All(ctx, &institutions); err != nil {
		return nil, fmt.Errorf("failed to decode financial institutions: %w", err)
	}

	return institutions, nil
}

// GetPendingSyncAdditionalBalances fetches additional balances that need to be synced
func (c *Connection) GetPendingSyncAdditionalBalances(ctx context.Context, limit int) ([]AdditionalBalanceDocument, error) {
	collection := c.Collection("additional_balances")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find additional balances: %w", err)
	}
	defer cursor.Close(ctx)

	var balances []AdditionalBalanceDocument
	if err := cursor.All(ctx, &balances); err != nil {
		return nil, fmt.Errorf("failed to decode additional balances: %w", err)
	}

	return balances, nil
}

// GetPendingSyncBalanceHistory fetches balance history records that need to be synced
func (c *Connection) GetPendingSyncBalanceHistory(ctx context.Context, limit int) ([]BalanceHistoryDocument, error) {
	collection := c.Collection("balance_history")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find balance history: %w", err)
	}
	defer cursor.Close(ctx)

	var history []BalanceHistoryDocument
	if err := cursor.All(ctx, &history); err != nil {
		return nil, fmt.Errorf("failed to decode balance history: %w", err)
	}

	return history, nil
}

// GetPendingSyncExpenseAutomaticWorkflows fetches expense automatic workflows that need to be synced
func (c *Connection) GetPendingSyncExpenseAutomaticWorkflows(ctx context.Context, limit int) ([]ExpenseAutomaticWorkflowDocument, error) {
	collection := c.Collection("expense_automatic_workflow")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find expense automatic workflows: %w", err)
	}
	defer cursor.Close(ctx)

	var workflows []ExpenseAutomaticWorkflowDocument
	if err := cursor.All(ctx, &workflows); err != nil {
		return nil, fmt.Errorf("failed to decode expense automatic workflows: %w", err)
	}

	return workflows, nil
}

// GetPendingSyncExpenseAutomaticWorkflowPreSavedDescriptions fetches pre-saved descriptions that need to be synced
func (c *Connection) GetPendingSyncExpenseAutomaticWorkflowPreSavedDescriptions(ctx context.Context, limit int) ([]ExpenseAutomaticWorkflowPreSavedDescriptionDocument, error) {
	collection := c.Collection("expense_automatic_workflow_pre_saved_description")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find pre-saved descriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var descriptions []ExpenseAutomaticWorkflowPreSavedDescriptionDocument
	if err := cursor.All(ctx, &descriptions); err != nil {
		return nil, fmt.Errorf("failed to decode pre-saved descriptions: %w", err)
	}

	return descriptions, nil
}

// GetPendingSyncServicePayments fetches service payments that need to be synced
func (c *Connection) GetPendingSyncServicePayments(ctx context.Context, limit int) ([]ServicePaymentDocument, error) {
	collection := c.Collection("payments")

	filter := bson.M{
		"$or": []bson.M{
			{"onPremiseSyncDatetime": bson.M{"$exists": false}},
			{"onPremiseSyncDatetime": nil},
		},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find service payments: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []ServicePaymentDocument
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode service payments: %w", err)
	}

	return payments, nil
}

// GetSettings fetches system settings
func (c *Connection) GetSettings(ctx context.Context) (*SettingsDocument, error) {
	collection := c.Collection("settings")

	var settings SettingsDocument
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}

	return &settings, nil
}

// MarkAsSynced marks a document as synced in MongoDB
// Uses fields: onPremiseRelationalDBSyncDatetime, onPremiseRelationalDBSyncService
func (c *Connection) MarkAsSynced(ctx context.Context, collectionName string, docID string, serviceName string) error {
	coll := c.Collection(collectionName)

	filter := bson.M{"_id": docID}
	update := bson.M{
		"$set": bson.M{
			"onPremiseRelationalDBSyncDatetime": time.Now(),
			"onPremiseRelationalDBSyncService":  serviceName,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark document as synced: %w", err)
	}

	return nil
}

// GetSuccessfullyIngestedFirestoreDoc fetches a document from succesfully_ingested_firestore_docs collection by ID
func (c *Connection) GetSuccessfullyIngestedFirestoreDoc(ctx context.Context, docID string) (*SuccessfullyIngestedFirestoreDocsDocument, error) {
	collection := c.Collection("succesfully_ingested_firestore_docs")

	// Try to convert docID to ObjectID first, then fall back to string
	var filter bson.M
	objectID, err := primitive.ObjectIDFromHex(docID)
	if err == nil {
		filter = bson.M{"_id": objectID}
		log.Printf("Querying succesfully_ingested_firestore_docs with ObjectID: %s", objectID.Hex())
	} else {
		filter = bson.M{"_id": docID}
		log.Printf("Querying succesfully_ingested_firestore_docs with string ID: %s", docID)
	}

	var doc SuccessfullyIngestedFirestoreDocsDocument
	err = collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Document not found in succesfully_ingested_firestore_docs for ID: %s", docID)
			return nil, fmt.Errorf("document not found for ID: %s", docID)
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	log.Printf("Found document with %d collections in map_collection_to_docs", len(doc.MapCollectionToDocs))
	return &doc, nil
}

// GetUsersByIDs fetches users by their IDs
func (c *Connection) GetUsersByIDs(ctx context.Context, ids []string) ([]UserDocument, error) {
	collection := c.Collection("users")

	log.Printf("Fetching %d users from MongoDB", len(ids))
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []UserDocument
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	log.Printf("Retrieved %d users from MongoDB (requested %d)", len(users), len(ids))
	if len(users) != len(ids) {
		log.Printf("Warning: Some user IDs were not found in MongoDB")
	}

	return users, nil
}

// GetExpensesByIDs fetches expenses by their IDs
func (c *Connection) GetExpensesByIDs(ctx context.Context, ids []string) ([]ExpenseDocument, error) {
	collection := c.Collection("expenses")

	log.Printf("Fetching %d expenses from MongoDB", len(ids))
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []ExpenseDocument
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, fmt.Errorf("failed to decode expenses: %w", err)
	}

	log.Printf("Retrieved %d expenses from MongoDB (requested %d)", len(expenses), len(ids))
	if len(expenses) != len(ids) {
		log.Printf("Warning: Some expense IDs were not found in MongoDB")
	}

	return expenses, nil
}

// GetFinancialInstitutionsByIDs fetches financial institutions by their IDs
func (c *Connection) GetFinancialInstitutionsByIDs(ctx context.Context, ids []string) ([]FinancialInstitutionDocument, error) {
	collection := c.Collection("banks")

	log.Printf("Fetching %d financial institutions from MongoDB", len(ids))
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find financial institutions: %w", err)
	}
	defer cursor.Close(ctx)

	var institutions []FinancialInstitutionDocument
	if err := cursor.All(ctx, &institutions); err != nil {
		return nil, fmt.Errorf("failed to decode financial institutions: %w", err)
	}

	log.Printf("Retrieved %d financial institutions from MongoDB (requested %d)", len(institutions), len(ids))

	return institutions, nil
}

// GetAdditionalBalancesByIDs fetches additional balances by their IDs
func (c *Connection) GetAdditionalBalancesByIDs(ctx context.Context, ids []string) ([]AdditionalBalanceDocument, error) {
	collection := c.Collection("additional_balances")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find additional balances: %w", err)
	}
	defer cursor.Close(ctx)

	var balances []AdditionalBalanceDocument
	if err := cursor.All(ctx, &balances); err != nil {
		return nil, fmt.Errorf("failed to decode additional balances: %w", err)
	}

	return balances, nil
}

// GetBalanceHistoryByIDs fetches balance history records by their IDs
func (c *Connection) GetBalanceHistoryByIDs(ctx context.Context, ids []string) ([]BalanceHistoryDocument, error) {
	collection := c.Collection("balance_history")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find balance history: %w", err)
	}
	defer cursor.Close(ctx)

	var history []BalanceHistoryDocument
	if err := cursor.All(ctx, &history); err != nil {
		return nil, fmt.Errorf("failed to decode balance history: %w", err)
	}

	return history, nil
}

// GetExpenseAutomaticWorkflowsByIDs fetches expense automatic workflows by their IDs
func (c *Connection) GetExpenseAutomaticWorkflowsByIDs(ctx context.Context, ids []string) ([]ExpenseAutomaticWorkflowDocument, error) {
	collection := c.Collection("expense_automatic_workflow")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find expense automatic workflows: %w", err)
	}
	defer cursor.Close(ctx)

	var workflows []ExpenseAutomaticWorkflowDocument
	if err := cursor.All(ctx, &workflows); err != nil {
		return nil, fmt.Errorf("failed to decode expense automatic workflows: %w", err)
	}

	return workflows, nil
}

// GetExpenseAutomaticWorkflowPreSavedDescriptionsByIDs fetches pre-saved descriptions by their IDs
func (c *Connection) GetExpenseAutomaticWorkflowPreSavedDescriptionsByIDs(ctx context.Context, ids []string) ([]ExpenseAutomaticWorkflowPreSavedDescriptionDocument, error) {
	collection := c.Collection("expense_automatic_workflow_pre_saved_description")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find pre-saved descriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var descriptions []ExpenseAutomaticWorkflowPreSavedDescriptionDocument
	if err := cursor.All(ctx, &descriptions); err != nil {
		return nil, fmt.Errorf("failed to decode pre-saved descriptions: %w", err)
	}

	return descriptions, nil
}

// GetServicePaymentsByIDs fetches service payments by their IDs
func (c *Connection) GetServicePaymentsByIDs(ctx context.Context, ids []string) ([]ServicePaymentDocument, error) {
	collection := c.Collection("payments")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find service payments: %w", err)
	}
	defer cursor.Close(ctx)

	var payments []ServicePaymentDocument
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode service payments: %w", err)
	}

	return payments, nil
}

// GetSettingsByIDs fetches settings documents by their IDs
func (c *Connection) GetSettingsByIDs(ctx context.Context, ids []string) ([]SettingsDocument, error) {
	collection := c.Collection("settings")

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}
	defer cursor.Close(ctx)

	var settings []SettingsDocument
	if err := cursor.All(ctx, &settings); err != nil {
		return nil, fmt.Errorf("failed to decode settings: %w", err)
	}

	return settings, nil
}

// GetExpenseAggregate fetches all expenses with the same name and validity for a specific user.
// This is used for invoice/savings aggregation where multiple MongoDB expense records
// with the same expense name and validity should be consolidated into a single expense
// record in MariaDB with multiple installments.
func (c *Connection) GetExpenseAggregate(ctx context.Context, userID string, expenseName string, validity string) ([]ExpenseDocument, error) {
	collection := c.Collection("expenses")

	filter := bson.M{
		"user":        userID,
		"expenseName": expenseName,
		"validity":    validity,
	}

	// Sort by spendingDate to ensure chronological order
	opts := options.Find().SetSort(bson.M{"spendingDate": 1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find expense aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []ExpenseDocument
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, fmt.Errorf("failed to decode expense aggregate: %w", err)
	}

	return expenses, nil
}
