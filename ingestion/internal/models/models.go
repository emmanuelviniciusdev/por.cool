package models

import (
	"database/sql"
	"time"
)

// Domain represents the domain table
type Domain struct {
	ID        int64          `json:"id"`
	GUID      string         `json:"guid"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy sql.NullString `json:"created_by"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
}

// FinancialInstitution represents the financial_institution table
type FinancialInstitution struct {
	ID              int64          `json:"id"`
	GUID            string         `json:"guid"`
	SourceID        string         `json:"source_id"`
	UserID          int64          `json:"user_id"`
	Name            string         `json:"name"`
	FlCreditCard    bool           `json:"fl_credit_card"`
	FlMoneyMovement bool           `json:"fl_money_movement"`
	FlInvestment    bool           `json:"fl_investment"`
	CreatedAt       time.Time      `json:"created_at"`
	CreatedBy       sql.NullString `json:"created_by"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	UpdatedBy       sql.NullString `json:"updated_by"`
}

// ExpenseAutomaticWorkflow represents the expense_automatic_workflow table
type ExpenseAutomaticWorkflow struct {
	ID                               int64          `json:"id"`
	GUID                             string         `json:"guid"`
	SourceID                         string         `json:"source_id"`
	UserID                           int64          `json:"user_id"`
	Base64Image                      sql.NullString `json:"base64_image"`
	Description                      sql.NullString `json:"description"`
	ExtractedExpenseContentFromImage sql.NullString `json:"extracted_expense_content_from_image"`
	SpendingDateYYYYMM               sql.NullString `json:"spending_date__YYYY_MM"`
	SyncProcessedDate                sql.NullTime   `json:"sync_processed_date"`
	IDSyncStatus                     sql.NullInt64  `json:"id_sync_status"`
	ProcessingMessage                sql.NullString `json:"processing_message"`
	CreatedAt                        time.Time      `json:"created_at"`
	CreatedBy                        sql.NullString `json:"created_by"`
	UpdatedAt                        sql.NullTime   `json:"updated_at"`
	UpdatedBy                        sql.NullString `json:"updated_by"`
}

// ExpenseAutomaticWorkflowPreSavedDescription represents the expense_automatic_workflow_pre_saved_description table
type ExpenseAutomaticWorkflowPreSavedDescription struct {
	ID          int64          `json:"id"`
	GUID        string         `json:"guid"`
	SourceID    string         `json:"source_id"`
	UserID      int64          `json:"user_id"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	CreatedBy   sql.NullString `json:"created_by"`
	UpdatedAt   sql.NullTime   `json:"updated_at"`
	UpdatedBy   sql.NullString `json:"updated_by"`
}

// Expense represents the expense table
type Expense struct {
	ID                                int64          `json:"id"`
	GUID                              string         `json:"guid"`
	SourceID                          string         `json:"source_id"`
	UserID                            int64          `json:"user_id"`
	SpendingDateYYYYMM                string         `json:"spending_date__YYYY_MM"`
	IDStatus                          sql.NullInt64  `json:"id_status"`
	IDType                            sql.NullInt64  `json:"id_type"`
	ValidityPeriodDate                sql.NullTime   `json:"validity_period_date"`
	FlIndeterminateValidityPeriodDate bool           `json:"fl_indeterminate_validity_period_date"`
	Name                              string         `json:"name"`
	TotalAmount                       float64        `json:"total_amount"`
	TotalPaidAmount                   float64        `json:"total_paid_amount"`
	CreatedAt                         time.Time      `json:"created_at"`
	CreatedBy                         sql.NullString `json:"created_by"`
	UpdatedAt                         sql.NullTime   `json:"updated_at"`
	UpdatedBy                         sql.NullString `json:"updated_by"`
}

// ExpenseInstallment represents the expense_installment table
type ExpenseInstallment struct {
	ID         int64          `json:"id"`
	GUID       string         `json:"guid"`
	ExpenseID  int64          `json:"expense_id"`
	Amount     float64        `json:"amount"`
	PaidAmount float64        `json:"paid_amount"`
	IDStatus   sql.NullInt64  `json:"id_status"`
	DueDate    sql.NullTime   `json:"due_date"`
	CreatedAt  time.Time      `json:"created_at"`
	CreatedBy  sql.NullString `json:"created_by"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
	UpdatedBy  sql.NullString `json:"updated_by"`
}

// AdditionalBalance represents the additional_balance table
type AdditionalBalance struct {
	ID                 int64          `json:"id"`
	GUID               string         `json:"guid"`
	SourceID           string         `json:"source_id"`
	UserID             int64          `json:"user_id"`
	SpendingDateYYYYMM string         `json:"spending_date__YYYY_MM"`
	Amount             float64        `json:"amount"`
	Description        sql.NullString `json:"description"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedBy          sql.NullString `json:"created_by"`
	UpdatedAt          sql.NullTime   `json:"updated_at"`
	UpdatedBy          sql.NullString `json:"updated_by"`
}

// BalanceHistory represents the balance_history table
type BalanceHistory struct {
	ID                 int64          `json:"id"`
	GUID               string         `json:"guid"`
	SourceID           string         `json:"source_id"`
	UserID             int64          `json:"user_id"`
	SpendingDateYYYYMM string         `json:"spending_date__YYYY_MM"`
	Amount             float64        `json:"amount"`
	LastMonthAmount    float64        `json:"last_month_amount"`
	MonthlyIncome      float64        `json:"monthly_income"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedBy          sql.NullString `json:"created_by"`
	UpdatedAt          sql.NullTime   `json:"updated_at"`
	UpdatedBy          sql.NullString `json:"updated_by"`
}

// ServicePayment represents the service_payment table
type ServicePayment struct {
	ID                   int64          `json:"id"`
	GUID                 string         `json:"guid"`
	SourceID             string         `json:"source_id"`
	UserID               int64          `json:"user_id"`
	ServicePaymentDate   time.Time      `json:"service_payment_date"`
	ServicePaymentTypeID sql.NullInt64  `json:"service_payment_type_id"`
	CreatedAt            time.Time      `json:"created_at"`
	CreatedBy            sql.NullString `json:"created_by"`
	UpdatedAt            sql.NullTime   `json:"updated_at"`
	UpdatedBy            sql.NullString `json:"updated_by"`
}

// User represents the user table
type User struct {
	ID                  int64          `json:"id"`
	GUID                string         `json:"guid"`
	SourceID            string         `json:"source_id"`
	FirstName           string         `json:"first_name"`
	LastName            sql.NullString `json:"last_name"`
	Email               string         `json:"email"`
	FlAdmin             bool           `json:"fl_admin"`
	MonthlyIncome       float64        `json:"monthly_income"`
	FlPaymentRequested  bool           `json:"fl_payment_requested"`
	FlPaymentPending    bool           `json:"fl_payment_pending"`
	FlPaymentPaid       bool           `json:"fl_payment_paid"`
	CurrentSpendingDate sql.NullString `json:"current_spending_date"`
	CreatedAt           time.Time      `json:"created_at"`
	CreatedBy           sql.NullString `json:"created_by"`
	UpdatedAt           sql.NullTime   `json:"updated_at"`
	UpdatedBy           sql.NullString `json:"updated_by"`
}

// SystemSettings represents the system_settings table
type SystemSettings struct {
	ID                      int64          `json:"id"`
	GUID                    string         `json:"guid"`
	SourceID                string         `json:"source_id"`
	FlBlockUserRegistration bool           `json:"fl_block_user_registration"`
	FlMaintenance           bool           `json:"fl_maintenance"`
	JSONSyncMetadata        sql.NullString `json:"json_sync_metadata"`
	CreatedAt               time.Time      `json:"created_at"`
	CreatedBy               sql.NullString `json:"created_by"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	UpdatedBy               sql.NullString `json:"updated_by"`
}

// DomainSeed represents a domain to be seeded
type DomainSeed struct {
	Source string
	Type   string
	Names  []string
}

// GetDomainSeeds returns all domains to be seeded
func GetDomainSeeds() []DomainSeed {
	return []DomainSeed{
		{
			Source: "expense_automatic_workflow",
			Type:   "id_sync_status",
			Names:  []string{"pending", "success", "error"},
		},
		{
			Source: "expense",
			Type:   "id_status",
			Names:  []string{"pending", "partially_paid", "paid"},
		},
		{
			Source: "expense",
			Type:   "id_type",
			Names:  []string{"expense", "invoice", "savings"},
		},
		{
			Source: "expense_installment",
			Type:   "id_status",
			Names:  []string{"pending", "partially_paid", "paid"},
		},
		{
			Source: "service_payment",
			Type:   "service_payment_type_id",
			Names:  []string{"PayPal"},
		},
	}
}
