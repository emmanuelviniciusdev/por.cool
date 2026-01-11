package ingestion

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/porcool/ingestion/internal/config"
	"github.com/porcool/ingestion/internal/database/mariadb"
	"github.com/porcool/ingestion/internal/database/mongodb"
	"github.com/porcool/ingestion/internal/models"
)

const serviceName = "porcool-ingestion-non-relational-database-to-relational-database"

// Service handles the ingestion process from MongoDB to MariaDB
type Service struct {
	mariaDB *mariadb.Connection
	mongoDB *mongodb.Connection
	cfg     *config.Config
}

// NewService creates a new ingestion service
func NewService(mariaDB *mariadb.Connection, mongoDB *mongodb.Connection, cfg *config.Config) *Service {
	return &Service{
		mariaDB: mariaDB,
		mongoDB: mongoDB,
		cfg:     cfg,
	}
}

// formatSpendingDate converts a spending date to the YYYY/MM format for MariaDB
// Input formats supported: "YYYY-MM", "YYYY/MM", "YYYYMM"
func formatSpendingDate(date string) string {
	if date == "" {
		return ""
	}

	// Remove any dashes or slashes and normalize
	normalized := strings.ReplaceAll(date, "-", "")
	normalized = strings.ReplaceAll(normalized, "/", "")

	// If it's already in YYYYMM format (6 chars), convert to YYYY/MM
	if len(normalized) == 6 {
		return normalized[:4] + "/" + normalized[4:]
	}

	// If the date is in YYYY-MM or YYYY/MM format, extract and reformat
	if len(date) >= 7 {
		year := date[:4]
		month := date[5:7]
		return year + "/" + month
	}

	// Return as-is if we can't parse it
	return date
}

// parseSpendingDateToTime parses a spending date string (YYYY/MM or YYYY-MM) to a time.Time
// The day is set to 1 (first day of the month)
func parseSpendingDateToTime(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}

	normalized := strings.ReplaceAll(date, "-", "/")
	parts := strings.Split(normalized, "/")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format: %s", date)
	}

	dateStr := fmt.Sprintf("%s-%s-01", parts[0], parts[1])
	return time.Parse("2006-01-02", dateStr)
}

// addMonths adds n months to a date and returns YYYY/MM format
func addMonths(date string, months int) string {
	t, err := parseSpendingDateToTime(date)
	if err != nil {
		return ""
	}
	t = t.AddDate(0, months, 0)
	return t.Format("2006/01")
}

// generateMonthRange generates all months between start and end (inclusive) in YYYY/MM format
func generateMonthRange(start, end string) []string {
	startTime, err := parseSpendingDateToTime(start)
	if err != nil {
		return nil
	}
	endTime, err := parseSpendingDateToTime(end)
	if err != nil {
		return nil
	}

	var months []string
	for t := startTime; !t.After(endTime); t = t.AddDate(0, 1, 0) {
		months = append(months, t.Format("2006/01"))
	}
	return months
}

// syncSimpleExpense creates a single expense record without installments
func (s *Service) syncSimpleExpense(ctx context.Context, mongoExpense mongodb.ExpenseDocument, userID int64) error {
	expenseRepo := mariadb.NewExpenseRepository(s.mariaDB)

	var statusID sql.NullInt64
	if mongoExpense.Status != "" {
		id, err := s.mariaDB.GetDomainID(mongoExpense.Status, "id_status", "expense")
		if err == nil {
			statusID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	var typeID sql.NullInt64
	if mongoExpense.Type != "" {
		id, err := s.mariaDB.GetDomainID(mongoExpense.Type, "id_type", "expense")
		if err == nil {
			typeID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	var validityDate sql.NullTime
	if mongoExpense.Validity != nil && *mongoExpense.Validity != "" {
		t, err := parseSpendingDateToTime(*mongoExpense.Validity)
		if err == nil {
			validityDate = sql.NullTime{Time: t, Valid: true}
		}
	}

	expense := &models.Expense{
		SourceID:                          mongoExpense.ID,
		UserID:                            userID,
		SpendingDateYYYYMM:                formatSpendingDate(mongoExpense.SpendingDate),
		IDStatus:                          statusID,
		IDType:                            typeID,
		ValidityPeriodDate:                validityDate,
		FlIndeterminateValidityPeriodDate: mongoExpense.IndeterminateValidity,
		Name:                              mongoExpense.ExpenseName,
		TotalAmount:                       mongoExpense.Amount,
		TotalPaidAmount:                   mongoExpense.AlreadyPaidAmount,
	}

	return expenseRepo.UpsertExpense(expense)
}

// syncExpenseWithInstallments handles invoice/savings with validity dates
// This fetches all related expenses (same name and validity) and generates installments
func (s *Service) syncExpenseWithInstallments(ctx context.Context, mongoExpense mongodb.ExpenseDocument, userID int64) error {
	expenseRepo := mariadb.NewExpenseRepository(s.mariaDB)
	installmentRepo := mariadb.NewExpenseInstallmentRepository(s.mariaDB)

	validity := *mongoExpense.Validity
	expenseName := mongoExpense.ExpenseName

	// Get all expenses in the aggregate (same name and validity for this user)
	aggregateExpenses, err := s.mongoDB.GetExpenseAggregate(ctx, mongoExpense.User, expenseName, validity)
	if err != nil {
		return fmt.Errorf("failed to get expense aggregate: %w", err)
	}

	// Format validity to YYYY/MM for database lookup
	// The validity from MongoDB can be in ISO format (2026-03-01T03:00:00Z) or YYYY-MM format
	validityFormatted := formatSpendingDate(validity)

	// If this aggregate has already been processed, just mark as synced and return
	// We check by looking for an existing expense with this name, validity, and no spending_date
	existingExpense, err := expenseRepo.GetExpenseByNameValidityUser(expenseName, validityFormatted, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing expense: %w", err)
	}

	var typeID sql.NullInt64
	if mongoExpense.Type != "" {
		id, err := s.mariaDB.GetDomainID(mongoExpense.Type, "id_type", "expense")
		if err == nil {
			typeID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	var validityDate sql.NullTime
	// Use the formatted validity (YYYY/MM) for parsing to time.Time
	t, err := parseSpendingDateToTime(validityFormatted)
	if err == nil {
		validityDate = sql.NullTime{Time: t, Valid: true}
	}

	var expenseID int64

	if existingExpense != nil {
		// Expense already exists, use its ID
		expenseID = existingExpense.ID
		log.Printf("Using existing expense record for aggregate: %s (ID: %d)", expenseName, expenseID)
	} else {
		// Create a "generic" expense record (no spending_date, id_status, total_amount, total_paid_amount)
		expense := &models.Expense{
			SourceID:                          mongoExpense.ID,
			UserID:                            userID,
			SpendingDateYYYYMM:                "", // Empty for aggregate expense
			IDStatus:                          sql.NullInt64{Valid: false},
			IDType:                            typeID,
			ValidityPeriodDate:                validityDate,
			FlIndeterminateValidityPeriodDate: mongoExpense.IndeterminateValidity,
			Name:                              expenseName,
			TotalAmount:                       0, // Not set for aggregate
			TotalPaidAmount:                   0, // Not set for aggregate
		}

		if err := expenseRepo.UpsertExpense(expense); err != nil {
			return fmt.Errorf("failed to create generic expense record: %w", err)
		}
		expenseID = expense.ID
		log.Printf("Created new generic expense record for aggregate: %s (ID: %d)", expenseName, expenseID)
	}

	// Sort aggregate expenses by spending date
	sort.Slice(aggregateExpenses, func(i, j int) bool {
		return aggregateExpenses[i].SpendingDate < aggregateExpenses[j].SpendingDate
	})

	// Create installments from MongoDB expense records
	existingInstallmentDates := make(map[string]bool)

	for _, aggExp := range aggregateExpenses {
		spendingDate := formatSpendingDate(aggExp.SpendingDate)

		// Check if installment already exists for this date
		existingInstallment, err := installmentRepo.GetInstallmentByExpenseAndDate(expenseID, spendingDate)
		if err != nil {
			log.Printf("Error checking existing installment: %v", err)
		}

		if existingInstallment != nil {
			// Update existing installment
			existingInstallment.Amount = aggExp.Amount
			existingInstallment.PaidAmount = aggExp.AlreadyPaidAmount

			var statusID sql.NullInt64
			if aggExp.Status != "" {
				id, err := s.mariaDB.GetDomainID(aggExp.Status, "id_status", "expense_installment")
				if err == nil {
					statusID = sql.NullInt64{Int64: id, Valid: true}
				}
			}
			existingInstallment.IDStatus = statusID

			dueDate, _ := parseSpendingDateToTime(spendingDate)
			existingInstallment.DueDate = sql.NullTime{Time: dueDate, Valid: true}

			if err := installmentRepo.UpsertExpenseInstallment(existingInstallment); err != nil {
				log.Printf("Error updating installment for date %s: %v", spendingDate, err)
			}
		} else {
			// Create new installment
			var statusID sql.NullInt64
			if aggExp.Status != "" {
				id, err := s.mariaDB.GetDomainID(aggExp.Status, "id_status", "expense_installment")
				if err == nil {
					statusID = sql.NullInt64{Int64: id, Valid: true}
				}
			}

			dueDate, _ := parseSpendingDateToTime(spendingDate)

			installment := &models.ExpenseInstallment{
				ExpenseID:  expenseID,
				Amount:     aggExp.Amount,
				PaidAmount: aggExp.AlreadyPaidAmount,
				IDStatus:   statusID,
				DueDate:    sql.NullTime{Time: dueDate, Valid: true},
			}

			if err := installmentRepo.UpsertExpenseInstallment(installment); err != nil {
				log.Printf("Error creating installment for date %s: %v", spendingDate, err)
			}
		}

		existingInstallmentDates[spendingDate] = true
	}

	// Generate remaining installments from the last MongoDB expense date + 1 month until validity
	if len(aggregateExpenses) > 0 {
		lastExpenseDate := formatSpendingDate(aggregateExpenses[len(aggregateExpenses)-1].SpendingDate)
		validityFormatted := formatSpendingDate(validity)

		// Generate months from lastExpenseDate + 1 month until validity
		nextMonth := addMonths(lastExpenseDate, 1)
		remainingMonths := generateMonthRange(nextMonth, validityFormatted)

		// Get the "pending" status ID for new installments
		pendingStatusID, _ := s.mariaDB.GetDomainID("pending", "id_status", "expense_installment")

		for _, month := range remainingMonths {
			if existingInstallmentDates[month] {
				continue // Skip if already exists
			}

			// Check if installment already exists in DB
			existingInstallment, _ := installmentRepo.GetInstallmentByExpenseAndDate(expenseID, month)
			if existingInstallment != nil {
				continue // Skip if already exists in DB
			}

			dueDate, _ := parseSpendingDateToTime(month)

			installment := &models.ExpenseInstallment{
				ExpenseID:  expenseID,
				Amount:     mongoExpense.Amount, // New generated installments use the MongoDB expense's amount
				PaidAmount: 0,
				IDStatus:   sql.NullInt64{Int64: pendingStatusID, Valid: true},
				DueDate:    sql.NullTime{Time: dueDate, Valid: true},
			}

			if err := installmentRepo.UpsertExpenseInstallment(installment); err != nil {
				log.Printf("Error creating generated installment for date %s: %v", month, err)
			} else {
				log.Printf("Generated pending installment for %s: %s", expenseName, month)
			}
		}
	}

	return nil
}

// ProcessIngestionMessage processes documents based on a successfully ingested firestore docs record
// This is the main entry point for message-based processing from RabbitMQ
func (s *Service) ProcessIngestionMessage(ctx context.Context, docID string) error {
	log.Printf("Processing ingestion message for document ID: %s", docID)

	// Fetch the document from succesfully_ingested_firestore_docs collection
	doc, err := s.mongoDB.GetSuccessfullyIngestedFirestoreDoc(ctx, docID)
	if err != nil {
		log.Printf("Error fetching document from MongoDB: %v", err)
		return fmt.Errorf("failed to fetch ingestion document: %w", err)
	}
	if doc == nil {
		log.Printf("Document not found for ID: %s - this should not happen", docID)
		return fmt.Errorf("document not found for ID: %s", docID)
	}

	log.Printf("Found document with %d collections to process", len(doc.MapCollectionToDocs))

	// Define the correct ingestion order - users must be first since other collections depend on them
	collectionOrder := []string{
		"users",
		"banks",
		"expenses",
		"additional_balances",
		"balance_history",
		"expense_automatic_workflow",
		"expense_automatic_workflow_pre_saved_description",
		"payments",
		"settings",
	}

	// Track errors for each collection
	var collectionErrors []string
	processedCollections := 0
	successfulCollections := 0

	// Process collections in the correct order
	for _, collectionName := range collectionOrder {
		docIDs, exists := doc.MapCollectionToDocs[collectionName]
		if !exists {
			continue
		}

		ids := extractDocIDs(docIDs)
		if len(ids) == 0 {
			log.Printf("No document IDs found for collection: %s", collectionName)
			continue
		}

		processedCollections++
		log.Printf("Processing %d documents from collection: %s", len(ids), collectionName)

		var syncErr error
		switch collectionName {
		case "users":
			syncErr = s.syncUsersByIDs(ctx, ids)
		case "expenses":
			syncErr = s.syncExpensesByIDs(ctx, ids)
		case "banks":
			syncErr = s.syncFinancialInstitutionsByIDs(ctx, ids)
		case "additional_balances":
			syncErr = s.syncAdditionalBalancesByIDs(ctx, ids)
		case "balance_history":
			syncErr = s.syncBalanceHistoryByIDs(ctx, ids)
		case "expense_automatic_workflow":
			syncErr = s.syncExpenseAutomaticWorkflowsByIDs(ctx, ids)
		case "expense_automatic_workflow_pre_saved_description":
			syncErr = s.syncExpenseAutomaticWorkflowPreSavedDescriptionsByIDs(ctx, ids)
		case "payments":
			syncErr = s.syncServicePaymentsByIDs(ctx, ids)
		case "settings":
			syncErr = s.syncSettingsByIDs(ctx, ids)
		default:
			log.Printf("Unknown collection: %s", collectionName)
			continue
		}

		if syncErr != nil {
			errMsg := fmt.Sprintf("%s: %v", collectionName, syncErr)
			log.Printf("Error syncing %s: %v", collectionName, syncErr)
			collectionErrors = append(collectionErrors, errMsg)
		} else {
			successfulCollections++
			log.Printf("Successfully synced collection: %s", collectionName)
		}
	}

	log.Printf("Completed processing ingestion message for document ID: %s (processed: %d, successful: %d, failed: %d)",
		docID, processedCollections, successfulCollections, len(collectionErrors))

	// Return error if any collection failed to sync
	if len(collectionErrors) > 0 {
		return fmt.Errorf("failed to sync %d collection(s): %s", len(collectionErrors), strings.Join(collectionErrors, "; "))
	}

	// Mark the ingestion document as processed with ingestedBy and ingestedAt
	if err := s.mongoDB.MarkIngestionDocAsProcessed(ctx, docID, serviceName); err != nil {
		log.Printf("Warning: failed to mark ingestion document as processed: %v", err)
		// Don't return error here - the ingestion was successful, this is just metadata
	}

	return nil
}

// extractDocIDs extracts string document IDs from various formats
// The map_collection_to_docs field can contain arrays of strings or other formats
func extractDocIDs(docIDs interface{}) []string {
	var ids []string

	log.Printf("extractDocIDs: received type %T, value: %v", docIDs, docIDs)

	switch v := docIDs.(type) {
	case []interface{}:
		for _, id := range v {
			log.Printf("extractDocIDs: array element type %T, value: %v", id, id)
			if strID, ok := id.(string); ok {
				ids = append(ids, strID)
			}
		}
	case []string:
		ids = v
	case string:
		ids = []string{v}
	case primitive.A:
		// MongoDB primitive.A is the BSON array type
		for _, id := range v {
			log.Printf("extractDocIDs: primitive.A element type %T, value: %v", id, id)
			if strID, ok := id.(string); ok {
				ids = append(ids, strID)
			}
		}
	default:
		log.Printf("extractDocIDs: unhandled type %T", docIDs)
	}

	log.Printf("extractDocIDs: extracted %d IDs", len(ids))
	return ids
}

// syncUsersByIDs syncs specific users by their IDs
func (s *Service) syncUsersByIDs(ctx context.Context, ids []string) error {
	users, err := s.mongoDB.GetUsersByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d users to sync", len(users))

	repo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoUser := range users {
		user := &models.User{
			SourceID:            mongoUser.ID,
			FirstName:           mongoUser.Name,
			LastName:            sql.NullString{String: mongoUser.LastName, Valid: mongoUser.LastName != ""},
			Email:               mongoUser.Email,
			FlAdmin:             mongoUser.Admin,
			MonthlyIncome:       mongoUser.MonthlyIncome,
			FlPaymentRequested:  mongoUser.RequestedPayment,
			FlPaymentPending:    mongoUser.PendingPayment,
			FlPaymentPaid:       mongoUser.PaidPayment,
			CurrentSpendingDate: sql.NullString{String: formatSpendingDate(mongoUser.LookingAtSpendingDate), Valid: mongoUser.LookingAtSpendingDate != ""},
		}

		if err := repo.UpsertUser(user); err != nil {
			log.Printf("Error upserting user %s: %v", mongoUser.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "users", mongoUser.ID, serviceName); err != nil {
			log.Printf("Error marking user %s as synced: %v", mongoUser.ID, err)
		}

		log.Printf("Synced user: %s (%s)", user.Email, user.GUID)
	}

	return nil
}

// syncExpensesByIDs syncs specific expenses by their IDs
func (s *Service) syncExpensesByIDs(ctx context.Context, ids []string) error {
	expenses, err := s.mongoDB.GetExpensesByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d expenses to sync", len(expenses))

	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoExpense := range expenses {
		user, err := userRepo.GetUserBySourceID(mongoExpense.User)
		if err != nil || user == nil {
			log.Printf("User not found for expense %s: %v", mongoExpense.ID, err)
			continue
		}

		expenseType := mongoExpense.Type

		if expenseType == "expense" {
			if err := s.syncSimpleExpense(ctx, mongoExpense, user.ID); err != nil {
				log.Printf("Error syncing simple expense %s: %v", mongoExpense.ID, err)
				continue
			}
		} else if expenseType == "invoice" || expenseType == "savings" {
			if mongoExpense.Validity == nil || *mongoExpense.Validity == "" {
				if err := s.syncSimpleExpense(ctx, mongoExpense, user.ID); err != nil {
					log.Printf("Error syncing expense (no validity) %s: %v", mongoExpense.ID, err)
					continue
				}
			} else {
				if err := s.syncExpenseWithInstallments(ctx, mongoExpense, user.ID); err != nil {
					log.Printf("Error syncing expense with installments %s: %v", mongoExpense.ID, err)
					continue
				}
			}
		} else {
			if err := s.syncSimpleExpense(ctx, mongoExpense, user.ID); err != nil {
				log.Printf("Error syncing expense (unknown type) %s: %v", mongoExpense.ID, err)
				continue
			}
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "expenses", mongoExpense.ID, serviceName); err != nil {
			log.Printf("Error marking expense %s as synced: %v", mongoExpense.ID, err)
		}

		log.Printf("Synced expense: %s (%s)", mongoExpense.ExpenseName, mongoExpense.ID)
	}

	return nil
}

// syncFinancialInstitutionsByIDs syncs specific financial institutions by their IDs
func (s *Service) syncFinancialInstitutionsByIDs(ctx context.Context, ids []string) error {
	institutions, err := s.mongoDB.GetFinancialInstitutionsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d financial institutions to sync", len(institutions))

	fiRepo := mariadb.NewFinancialInstitutionRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoFI := range institutions {
		user, err := userRepo.GetUserBySourceID(mongoFI.User)
		if err != nil || user == nil {
			log.Printf("User not found for financial institution %s: %v", mongoFI.ID, err)
			continue
		}

		fi := &models.FinancialInstitution{
			SourceID:        mongoFI.ID,
			UserID:          user.ID,
			Name:            mongoFI.Nome,
			FlCreditCard:    mongoFI.CartaoCredito,
			FlMoneyMovement: mongoFI.MovimentacaoDinheiro,
			FlInvestment:    mongoFI.Investimentos,
		}

		if err := fiRepo.UpsertFinancialInstitution(fi); err != nil {
			log.Printf("Error upserting financial institution %s: %v", mongoFI.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "banks", mongoFI.ID, serviceName); err != nil {
			log.Printf("Error marking financial institution %s as synced: %v", mongoFI.ID, err)
		}

		log.Printf("Synced financial institution: %s (%s)", fi.Name, fi.GUID)
	}

	return nil
}

// syncAdditionalBalancesByIDs syncs specific additional balances by their IDs
func (s *Service) syncAdditionalBalancesByIDs(ctx context.Context, ids []string) error {
	balances, err := s.mongoDB.GetAdditionalBalancesByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d additional balances to sync", len(balances))

	abRepo := mariadb.NewAdditionalBalanceRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoAB := range balances {
		user, err := userRepo.GetUserBySourceID(mongoAB.User)
		if err != nil || user == nil {
			log.Printf("User not found for additional balance %s: %v", mongoAB.ID, err)
			continue
		}

		ab := &models.AdditionalBalance{
			SourceID:           mongoAB.ID,
			UserID:             user.ID,
			SpendingDateYYYYMM: formatSpendingDate(mongoAB.SpendingDate),
			Amount:             mongoAB.Balance,
			Description:        sql.NullString{String: mongoAB.Description, Valid: mongoAB.Description != ""},
		}

		if err := abRepo.UpsertAdditionalBalance(ab); err != nil {
			log.Printf("Error upserting additional balance %s: %v", mongoAB.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "additional_balances", mongoAB.ID, serviceName); err != nil {
			log.Printf("Error marking additional balance %s as synced: %v", mongoAB.ID, err)
		}

		log.Printf("Synced additional balance: %s", ab.GUID)
	}

	return nil
}

// syncBalanceHistoryByIDs syncs specific balance history records by their IDs
func (s *Service) syncBalanceHistoryByIDs(ctx context.Context, ids []string) error {
	history, err := s.mongoDB.GetBalanceHistoryByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d balance history records to sync", len(history))

	bhRepo := mariadb.NewBalanceHistoryRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoBH := range history {
		user, err := userRepo.GetUserBySourceID(mongoBH.User)
		if err != nil || user == nil {
			log.Printf("User not found for balance history %s: %v", mongoBH.ID, err)
			continue
		}

		bh := &models.BalanceHistory{
			SourceID:           mongoBH.ID,
			UserID:             user.ID,
			SpendingDateYYYYMM: formatSpendingDate(mongoBH.SpendingDate),
			Amount:             mongoBH.Balance,
			LastMonthAmount:    mongoBH.LastMonthBalance,
			MonthlyIncome:      mongoBH.MonthlyIncome,
		}

		if err := bhRepo.UpsertBalanceHistory(bh); err != nil {
			log.Printf("Error upserting balance history %s: %v", mongoBH.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "balance_history", mongoBH.ID, serviceName); err != nil {
			log.Printf("Error marking balance history %s as synced: %v", mongoBH.ID, err)
		}

		log.Printf("Synced balance history: %s", bh.GUID)
	}

	return nil
}

// syncExpenseAutomaticWorkflowsByIDs syncs specific expense automatic workflows by their IDs
func (s *Service) syncExpenseAutomaticWorkflowsByIDs(ctx context.Context, ids []string) error {
	workflows, err := s.mongoDB.GetExpenseAutomaticWorkflowsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d expense automatic workflows to sync", len(workflows))

	eawRepo := mariadb.NewExpenseAutomaticWorkflowRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoEAW := range workflows {
		user, err := userRepo.GetUserBySourceID(mongoEAW.User)
		if err != nil || user == nil {
			log.Printf("User not found for expense automatic workflow %s: %v", mongoEAW.ID, err)
			continue
		}

		var syncStatusID sql.NullInt64
		if mongoEAW.SyncStatus != "" {
			id, err := s.mariaDB.GetDomainID(mongoEAW.SyncStatus, "id_sync_status", "expense_automatic_workflow")
			if err == nil {
				syncStatusID = sql.NullInt64{Int64: id, Valid: true}
			}
		}

		var extractedContent sql.NullString
		if mongoEAW.ExtractedExpenseContentFromImage != nil {
			jsonBytes, err := json.Marshal(mongoEAW.ExtractedExpenseContentFromImage)
			if err == nil && string(jsonBytes) != "null" {
				extractedContent = sql.NullString{String: string(jsonBytes), Valid: true}
			}
		}

		var syncProcessedDate sql.NullTime
		if mongoEAW.SyncProcessedDate != "" {
			t, err := time.Parse(time.RFC3339, mongoEAW.SyncProcessedDate)
			if err == nil {
				syncProcessedDate = sql.NullTime{Time: t, Valid: true}
			}
		}

		eaw := &models.ExpenseAutomaticWorkflow{
			SourceID:                         mongoEAW.ID,
			UserID:                           user.ID,
			Base64Image:                      sql.NullString{String: mongoEAW.Base64Image, Valid: mongoEAW.Base64Image != ""},
			Description:                      sql.NullString{String: mongoEAW.Description, Valid: mongoEAW.Description != ""},
			ExtractedExpenseContentFromImage: extractedContent,
			SpendingDateYYYYMM:               sql.NullString{String: formatSpendingDate(mongoEAW.SpendingDate), Valid: mongoEAW.SpendingDate != ""},
			SyncProcessedDate:                syncProcessedDate,
			IDSyncStatus:                     syncStatusID,
			ProcessingMessage:                sql.NullString{String: mongoEAW.ProcessingMessage, Valid: mongoEAW.ProcessingMessage != ""},
		}

		if err := eawRepo.UpsertExpenseAutomaticWorkflow(eaw); err != nil {
			log.Printf("Error upserting expense automatic workflow %s: %v", mongoEAW.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "expense_automatic_workflow", mongoEAW.ID, serviceName); err != nil {
			log.Printf("Error marking expense automatic workflow %s as synced: %v", mongoEAW.ID, err)
		}

		log.Printf("Synced expense automatic workflow: %s", eaw.GUID)
	}

	return nil
}

// syncExpenseAutomaticWorkflowPreSavedDescriptionsByIDs syncs specific pre-saved descriptions by their IDs
func (s *Service) syncExpenseAutomaticWorkflowPreSavedDescriptionsByIDs(ctx context.Context, ids []string) error {
	descriptions, err := s.mongoDB.GetExpenseAutomaticWorkflowPreSavedDescriptionsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d expense automatic workflow pre-saved descriptions to sync", len(descriptions))

	eawpsdRepo := mariadb.NewExpenseAutomaticWorkflowPreSavedDescriptionRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoDesc := range descriptions {
		user, err := userRepo.GetUserBySourceID(mongoDesc.User)
		if err != nil || user == nil {
			log.Printf("User not found for pre-saved description %s: %v", mongoDesc.ID, err)
			continue
		}

		desc := &models.ExpenseAutomaticWorkflowPreSavedDescription{
			SourceID:    mongoDesc.ID,
			UserID:      user.ID,
			Description: mongoDesc.Description,
		}

		if err := eawpsdRepo.UpsertExpenseAutomaticWorkflowPreSavedDescription(desc); err != nil {
			log.Printf("Error upserting pre-saved description %s: %v", mongoDesc.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "expense_automatic_workflow_pre_saved_description", mongoDesc.ID, serviceName); err != nil {
			log.Printf("Error marking pre-saved description %s as synced: %v", mongoDesc.ID, err)
		}

		log.Printf("Synced pre-saved description: %s", desc.GUID)
	}

	return nil
}

// syncServicePaymentsByIDs syncs specific service payments by their IDs
func (s *Service) syncServicePaymentsByIDs(ctx context.Context, ids []string) error {
	payments, err := s.mongoDB.GetServicePaymentsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d service payments to sync", len(payments))

	spRepo := mariadb.NewServicePaymentRepository(s.mariaDB)
	userRepo := mariadb.NewUserRepository(s.mariaDB)

	for _, mongoSP := range payments {
		user, err := userRepo.GetUserBySourceID(mongoSP.User)
		if err != nil || user == nil {
			log.Printf("User not found for service payment %s: %v", mongoSP.ID, err)
			continue
		}

		var paymentDate time.Time
		if mongoSP.PaymentDate != "" {
			t, err := time.Parse("2006-01-02", mongoSP.PaymentDate)
			if err != nil {
				t, err = time.Parse(time.RFC3339, mongoSP.PaymentDate)
				if err != nil {
					log.Printf("Error parsing payment date for %s: %v", mongoSP.ID, err)
					continue
				}
			}
			paymentDate = t
		}

		paymentTypeID, err := s.mariaDB.GetDomainID("PayPal", "service_payment_type_id", "service_payment")
		var paymentTypeIDSQL sql.NullInt64
		if err == nil {
			paymentTypeIDSQL = sql.NullInt64{Int64: paymentTypeID, Valid: true}
		}

		sp := &models.ServicePayment{
			SourceID:             mongoSP.ID,
			UserID:               user.ID,
			ServicePaymentDate:   paymentDate,
			ServicePaymentTypeID: paymentTypeIDSQL,
		}

		if err := spRepo.UpsertServicePayment(sp); err != nil {
			log.Printf("Error upserting service payment %s: %v", mongoSP.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "payments", mongoSP.ID, serviceName); err != nil {
			log.Printf("Error marking service payment %s as synced: %v", mongoSP.ID, err)
		}

		log.Printf("Synced service payment: %s", sp.GUID)
	}

	return nil
}

// syncSettingsByIDs syncs specific system settings by their IDs
func (s *Service) syncSettingsByIDs(ctx context.Context, ids []string) error {
	settings, err := s.mongoDB.GetSettingsByIDs(ctx, ids)
	if err != nil {
		return err
	}

	log.Printf("Found %d system settings to sync", len(settings))

	ssRepo := mariadb.NewSystemSettingsRepository(s.mariaDB)

	for _, mongoSettings := range settings {
		// Serialize syncMetadata to JSON
		var syncMetadataJSON sql.NullString
		if mongoSettings.SyncMetadata != nil && len(mongoSettings.SyncMetadata) > 0 {
			jsonBytes, err := json.Marshal(mongoSettings.SyncMetadata)
			if err == nil && string(jsonBytes) != "null" {
				syncMetadataJSON = sql.NullString{String: string(jsonBytes), Valid: true}
			}
		}

		ss := &models.SystemSettings{
			SourceID:                mongoSettings.ID,
			FlBlockUserRegistration: mongoSettings.BlockUserRegistration,
			FlMaintenance:           mongoSettings.Maintenance,
			JSONSyncMetadata:        syncMetadataJSON,
		}

		if err := ssRepo.UpsertSystemSettings(ss); err != nil {
			log.Printf("Error upserting system settings %s: %v", mongoSettings.ID, err)
			continue
		}

		if err := s.mongoDB.MarkAsSynced(ctx, "settings", mongoSettings.ID, serviceName); err != nil {
			log.Printf("Error marking system settings %s as synced: %v", mongoSettings.ID, err)
		}

		log.Printf("Synced system settings: %s", ss.GUID)
	}

	return nil
}
