package mongodb

import (
	"testing"
	"time"

	"github.com/porcool/ingestion/internal/config"
)

func TestMongoDBConfig(t *testing.T) {
	cfg := config.MongoDBConfig{
		URI:      "mongodb://localhost:27017",
		Database: "testdb",
	}

	if cfg.URI != "mongodb://localhost:27017" {
		t.Errorf("URI = %s, want mongodb://localhost:27017", cfg.URI)
	}

	if cfg.Database != "testdb" {
		t.Errorf("Database = %s, want testdb", cfg.Database)
	}
}

func TestUserDocumentStructure(t *testing.T) {
	user := UserDocument{
		ID:                    "user-123",
		FirestoreCreateTime:   "2024-01-01T00:00:00Z",
		FirestorePath:         "users/user-123",
		FirestoreUpdateTime:   "2024-01-02T00:00:00Z",
		ImportedAt:            time.Now(),
		Admin:                 true,
		Email:                 "john@example.com",
		LastName:              "Doe",
		LookingAtSpendingDate: "2024-01",
		MonthlyIncome:         5000.00,
		Name:                  "John",
		PaidPayment:           true,
		RequestedPayment:      false,
		PendingPayment:        false,
	}

	if user.ID != "user-123" {
		t.Errorf("ID = %s, want user-123", user.ID)
	}
	if user.Name != "John" {
		t.Errorf("Name = %s, want John", user.Name)
	}
	if user.LastName != "Doe" {
		t.Errorf("LastName = %s, want Doe", user.LastName)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Email = %s, want john@example.com", user.Email)
	}
	if !user.Admin {
		t.Error("Admin = false, want true")
	}
	if user.MonthlyIncome != 5000.00 {
		t.Errorf("MonthlyIncome = %f, want 5000.00", user.MonthlyIncome)
	}
	if user.LookingAtSpendingDate != "2024-01" {
		t.Errorf("LookingAtSpendingDate = %s, want 2024-01", user.LookingAtSpendingDate)
	}
	if !user.PaidPayment {
		t.Error("PaidPayment = false, want true")
	}
}

func TestExpenseDocumentStructure(t *testing.T) {
	validity := "2024-12-31"
	expense := ExpenseDocument{
		ID:                    "expense-123",
		FirestoreCreateTime:   "2024-01-01T00:00:00Z",
		FirestorePath:         "expenses/expense-123",
		AlreadyPaidAmount:     50.25,
		Amount:                100.50,
		Created:               "2024-01-01T00:00:00Z",
		ExpenseName:           "Test Expense",
		IndeterminateValidity: false,
		Source:                "manual",
		SpendingDate:          "2024-01",
		Status:                "pending",
		Type:                  "expense",
		Updated:               "2024-01-02T00:00:00Z",
		User:                  "user-123",
		Validity:              &validity,
	}

	if expense.ID != "expense-123" {
		t.Errorf("ID = %s, want expense-123", expense.ID)
	}
	if expense.User != "user-123" {
		t.Errorf("User = %s, want user-123", expense.User)
	}
	if expense.SpendingDate != "2024-01" {
		t.Errorf("SpendingDate = %s, want 2024-01", expense.SpendingDate)
	}
	if expense.Status != "pending" {
		t.Errorf("Status = %s, want pending", expense.Status)
	}
	if expense.Amount != 100.50 {
		t.Errorf("Amount = %f, want 100.50", expense.Amount)
	}
	if expense.AlreadyPaidAmount != 50.25 {
		t.Errorf("AlreadyPaidAmount = %f, want 50.25", expense.AlreadyPaidAmount)
	}
	if expense.ExpenseName != "Test Expense" {
		t.Errorf("ExpenseName = %s, want Test Expense", expense.ExpenseName)
	}
	if *expense.Validity != "2024-12-31" {
		t.Errorf("Validity = %s, want 2024-12-31", *expense.Validity)
	}
}

func TestFinancialInstitutionDocumentStructure(t *testing.T) {
	fi := FinancialInstitutionDocument{
		ID:                   "fi-123",
		FirestoreCreateTime:  "2024-01-01T00:00:00Z",
		FirestorePath:        "banks/fi-123",
		CartaoCredito:        true,
		Created:              "2024-01-01T00:00:00Z",
		Investimentos:        false,
		MovimentacaoDinheiro: true,
		Nome:                 "Test Bank",
		Observacoes:          "Notes here",
		Updated:              "2024-01-02T00:00:00Z",
		User:                 "user-123",
	}

	if fi.ID != "fi-123" {
		t.Errorf("ID = %s, want fi-123", fi.ID)
	}
	if fi.Nome != "Test Bank" {
		t.Errorf("Nome = %s, want Test Bank", fi.Nome)
	}
	if !fi.CartaoCredito {
		t.Error("CartaoCredito = false, want true")
	}
	if !fi.MovimentacaoDinheiro {
		t.Error("MovimentacaoDinheiro = false, want true")
	}
	if fi.Investimentos {
		t.Error("Investimentos = true, want false")
	}
	if fi.User != "user-123" {
		t.Errorf("User = %s, want user-123", fi.User)
	}
}

func TestAdditionalBalanceDocumentStructure(t *testing.T) {
	ab := AdditionalBalanceDocument{
		ID:                  "ab-123",
		FirestoreCreateTime: "2024-01-01T00:00:00Z",
		FirestorePath:       "additional_balances/ab-123",
		Balance:             500.00,
		Created:             "2024-01-01T00:00:00Z",
		Description:         "Bonus",
		SpendingDate:        "2024-01",
		User:                "user-123",
	}

	if ab.ID != "ab-123" {
		t.Errorf("ID = %s, want ab-123", ab.ID)
	}
	if ab.Balance != 500.00 {
		t.Errorf("Balance = %f, want 500.00", ab.Balance)
	}
	if ab.Description != "Bonus" {
		t.Errorf("Description = %s, want Bonus", ab.Description)
	}
	if ab.SpendingDate != "2024-01" {
		t.Errorf("SpendingDate = %s, want 2024-01", ab.SpendingDate)
	}
	if ab.User != "user-123" {
		t.Errorf("User = %s, want user-123", ab.User)
	}
}

func TestBalanceHistoryDocumentStructure(t *testing.T) {
	bh := BalanceHistoryDocument{
		ID:                  "bh-123",
		FirestoreCreateTime: "2024-01-01T00:00:00Z",
		FirestorePath:       "balance_history/bh-123",
		Balance:             1000.00,
		Created:             "2024-01-01T00:00:00Z",
		LastMonthBalance:    800.00,
		MonthlyIncome:       5000.00,
		SpendingDate:        "2024-01",
		User:                "user-123",
	}

	if bh.ID != "bh-123" {
		t.Errorf("ID = %s, want bh-123", bh.ID)
	}
	if bh.Balance != 1000.00 {
		t.Errorf("Balance = %f, want 1000.00", bh.Balance)
	}
	if bh.LastMonthBalance != 800.00 {
		t.Errorf("LastMonthBalance = %f, want 800.00", bh.LastMonthBalance)
	}
	if bh.MonthlyIncome != 5000.00 {
		t.Errorf("MonthlyIncome = %f, want 5000.00", bh.MonthlyIncome)
	}
	if bh.User != "user-123" {
		t.Errorf("User = %s, want user-123", bh.User)
	}
}

func TestExpenseAutomaticWorkflowDocumentStructure(t *testing.T) {
	eaw := ExpenseAutomaticWorkflowDocument{
		ID:                  "eaw-123",
		FirestoreCreateTime: "2024-01-01T00:00:00Z",
		FirestorePath:       "expense_automatic_workflow/eaw-123",
		Base64Image:         "base64data",
		Created:             "2024-01-01T00:00:00Z",
		Description:         "Test workflow",
		SpendingDate:        "2024-01",
		SyncProcessedDate:   "2024-01-02T00:00:00Z",
		SyncStatus:          "pending",
		ProcessingMessage:   "",
		User:                "user-123",
	}

	if eaw.ID != "eaw-123" {
		t.Errorf("ID = %s, want eaw-123", eaw.ID)
	}
	if eaw.SyncStatus != "pending" {
		t.Errorf("SyncStatus = %s, want pending", eaw.SyncStatus)
	}
	if eaw.Base64Image != "base64data" {
		t.Errorf("Base64Image = %s, want base64data", eaw.Base64Image)
	}
	if eaw.User != "user-123" {
		t.Errorf("User = %s, want user-123", eaw.User)
	}
}

func TestExpenseAutomaticWorkflowPreSavedDescriptionDocumentStructure(t *testing.T) {
	desc := ExpenseAutomaticWorkflowPreSavedDescriptionDocument{
		ID:                  "desc-123",
		FirestoreCreateTime: "2024-01-01T00:00:00Z",
		FirestorePath:       "expense_automatic_workflow_pre_saved_description/desc-123",
		Created:             "2024-01-01T00:00:00Z",
		Description:         "Groceries",
		User:                "user-123",
	}

	if desc.ID != "desc-123" {
		t.Errorf("ID = %s, want desc-123", desc.ID)
	}
	if desc.Description != "Groceries" {
		t.Errorf("Description = %s, want Groceries", desc.Description)
	}
	if desc.User != "user-123" {
		t.Errorf("User = %s, want user-123", desc.User)
	}
}

func TestServicePaymentDocumentStructure(t *testing.T) {
	sp := ServicePaymentDocument{
		ID:                  "sp-123",
		FirestoreCreateTime: "2024-01-01T00:00:00Z",
		FirestorePath:       "payments/sp-123",
		PaymentDate:         "2024-01-15",
		User:                "user-123",
	}

	if sp.ID != "sp-123" {
		t.Errorf("ID = %s, want sp-123", sp.ID)
	}
	if sp.PaymentDate != "2024-01-15" {
		t.Errorf("PaymentDate = %s, want 2024-01-15", sp.PaymentDate)
	}
	if sp.User != "user-123" {
		t.Errorf("User = %s, want user-123", sp.User)
	}
}

func TestSettingsDocumentStructure(t *testing.T) {
	settings := SettingsDocument{
		ID:                    "settings-123",
		FirestoreCreateTime:   "2024-01-01T00:00:00Z",
		BlockUserRegistration: false,
		Maintenance:           false,
		SyncMetadata:          []interface{}{"item1", "item2"},
	}

	if settings.BlockUserRegistration {
		t.Error("BlockUserRegistration = true, want false")
	}
	if settings.Maintenance {
		t.Error("Maintenance = true, want false")
	}
	if len(settings.SyncMetadata) != 2 {
		t.Errorf("SyncMetadata length = %d, want 2", len(settings.SyncMetadata))
	}
}

func TestSuccessfullyIngestedFirestoreDocsDocumentStructure(t *testing.T) {
	doc := SuccessfullyIngestedFirestoreDocsDocument{
		ID:                   "doc-123",
		CreatedAt:            time.Now(),
		OnPremiseSyncService: "porcool-ingestion",
		MapCollectionToDocs:  map[string]interface{}{"users": []string{"user1", "user2"}},
	}

	if doc.OnPremiseSyncService != "porcool-ingestion" {
		t.Errorf("OnPremiseSyncService = %s, want porcool-ingestion", doc.OnPremiseSyncService)
	}
	if doc.MapCollectionToDocs["users"] == nil {
		t.Error("MapCollectionToDocs[users] should not be nil")
	}
}

func TestNewConnection_InvalidURI(t *testing.T) {
	cfg := config.MongoDBConfig{
		URI:      "invalid://uri",
		Database: "testdb",
	}

	_, err := NewConnection(cfg)
	if err == nil {
		t.Error("NewConnection() should fail with invalid URI")
	}
}
