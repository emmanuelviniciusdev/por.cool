package mariadb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/porcool/ingestion/internal/models"
)

// ServiceName is the name used for created_by/updated_by fields
const ServiceName = "porcool-ingestion-non-relational-database-to-relational-database"

// UserRepository handles user database operations
type UserRepository struct {
	conn *Connection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(conn *Connection) *UserRepository {
	return &UserRepository{conn: conn}
}

// UpsertUser inserts or updates a user
func (r *UserRepository) UpsertUser(user *models.User) error {
	// Check if user exists by GUID
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM user WHERE guid = ?", user.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new user
		result, err := r.conn.db.Exec(`
			INSERT INTO user (guid, source_id, first_name, last_name, email, fl_admin, monthly_income,
				fl_payment_requested, fl_payment_pending, fl_payment_paid, current_spending_date,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.GUID, user.SourceID, user.FirstName, user.LastName, user.Email, user.FlAdmin, user.MonthlyIncome,
			user.FlPaymentRequested, user.FlPaymentPending, user.FlPaymentPaid, user.CurrentSpendingDate,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
		id, _ := result.LastInsertId()
		user.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Update existing user
	_, err = r.conn.db.Exec(`
		UPDATE user SET source_id = ?, first_name = ?, last_name = ?, email = ?, fl_admin = ?, monthly_income = ?,
			fl_payment_requested = ?, fl_payment_pending = ?, fl_payment_paid = ?, current_spending_date = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		user.SourceID, user.FirstName, user.LastName, user.Email, user.FlAdmin, user.MonthlyIncome,
		user.FlPaymentRequested, user.FlPaymentPending, user.FlPaymentPaid, user.CurrentSpendingDate,
		time.Now(), ServiceName, user.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	user.ID = existingID
	return nil
}

// GetUserByGUID retrieves a user by GUID
func (r *UserRepository) GetUserByGUID(guid string) (*models.User, error) {
	user := &models.User{}
	err := r.conn.db.QueryRow(`
		SELECT id, guid, source_id, first_name, last_name, email, fl_admin, monthly_income,
			fl_payment_requested, fl_payment_pending, fl_payment_paid, current_spending_date,
			created_at, created_by, updated_at, updated_by
		FROM user WHERE guid = ?`, guid,
	).Scan(
		&user.ID, &user.GUID, &user.SourceID, &user.FirstName, &user.LastName, &user.Email, &user.FlAdmin, &user.MonthlyIncome,
		&user.FlPaymentRequested, &user.FlPaymentPending, &user.FlPaymentPaid, &user.CurrentSpendingDate,
		&user.CreatedAt, &user.CreatedBy, &user.UpdatedAt, &user.UpdatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ExpenseRepository handles expense database operations
type ExpenseRepository struct {
	conn *Connection
}

// NewExpenseRepository creates a new ExpenseRepository
func NewExpenseRepository(conn *Connection) *ExpenseRepository {
	return &ExpenseRepository{conn: conn}
}

// UpsertExpense inserts or updates an expense
func (r *ExpenseRepository) UpsertExpense(expense *models.Expense) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM expense WHERE guid = ?", expense.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO expense (guid, source_id, user_id, spending_date__YYYY_MM, id_status, id_type,
				validity_period_date, fl_indeterminate_validity_period_date, name, total_amount, total_paid_amount,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			expense.GUID, expense.SourceID, expense.UserID, expense.SpendingDateYYYYMM, expense.IDStatus, expense.IDType,
			expense.ValidityPeriodDate, expense.FlIndeterminateValidityPeriodDate, expense.Name, expense.TotalAmount, expense.TotalPaidAmount,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert expense: %w", err)
		}
		id, _ := result.LastInsertId()
		expense.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check expense existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE expense SET source_id = ?, user_id = ?, spending_date__YYYY_MM = ?, id_status = ?, id_type = ?,
			validity_period_date = ?, fl_indeterminate_validity_period_date = ?, name = ?, total_amount = ?, total_paid_amount = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		expense.SourceID, expense.UserID, expense.SpendingDateYYYYMM, expense.IDStatus, expense.IDType,
		expense.ValidityPeriodDate, expense.FlIndeterminateValidityPeriodDate, expense.Name, expense.TotalAmount, expense.TotalPaidAmount,
		time.Now(), ServiceName, expense.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}
	expense.ID = existingID
	return nil
}

// GetExpenseByNameValidityUser retrieves an expense by name, validity, and user ID.
// This is used to find aggregate expense records for invoice/savings with validity dates.
// For aggregate expenses, the spending_date is empty (NULL or empty string).
// IMPORTANT: This method is critical for preventing duplicate expense records
// when processing invoice/savings with the same name and validity for a user.
func (r *ExpenseRepository) GetExpenseByNameValidityUser(name string, validity string, userID int64) (*models.Expense, error) {
	expense := &models.Expense{}

	// For aggregate expenses, we look for records with empty spending_date
	// that match the name, validity (via validity_period_date), and user_id
	err := r.conn.db.QueryRow(`
		SELECT id, guid, source_id, user_id, spending_date__YYYY_MM, id_status, id_type,
			validity_period_date, fl_indeterminate_validity_period_date, name, total_amount, total_paid_amount,
			created_at, created_by, updated_at, updated_by
		FROM expense
		WHERE user_id = ? AND name = ? AND (spending_date__YYYY_MM = '' OR spending_date__YYYY_MM IS NULL)
			AND DATE_FORMAT(validity_period_date, '%Y/%m') = ?`,
		userID, name, validity,
	).Scan(
		&expense.ID, &expense.GUID, &expense.SourceID, &expense.UserID, &expense.SpendingDateYYYYMM, &expense.IDStatus, &expense.IDType,
		&expense.ValidityPeriodDate, &expense.FlIndeterminateValidityPeriodDate, &expense.Name, &expense.TotalAmount, &expense.TotalPaidAmount,
		&expense.CreatedAt, &expense.CreatedBy, &expense.UpdatedAt, &expense.UpdatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get expense by name/validity/user: %w", err)
	}
	return expense, nil
}

// ExpenseInstallmentRepository handles expense installment database operations
type ExpenseInstallmentRepository struct {
	conn *Connection
}

// NewExpenseInstallmentRepository creates a new ExpenseInstallmentRepository
func NewExpenseInstallmentRepository(conn *Connection) *ExpenseInstallmentRepository {
	return &ExpenseInstallmentRepository{conn: conn}
}

// UpsertExpenseInstallment inserts or updates an expense installment
func (r *ExpenseInstallmentRepository) UpsertExpenseInstallment(installment *models.ExpenseInstallment) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM expense_installment WHERE guid = ?", installment.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO expense_installment (guid, expense_id, amount, paid_amount, id_status, due_date,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			installment.GUID, installment.ExpenseID, installment.Amount, installment.PaidAmount, installment.IDStatus, installment.DueDate,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert expense installment: %w", err)
		}
		id, _ := result.LastInsertId()
		installment.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check expense installment existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE expense_installment SET expense_id = ?, amount = ?, paid_amount = ?, id_status = ?, due_date = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		installment.ExpenseID, installment.Amount, installment.PaidAmount, installment.IDStatus, installment.DueDate,
		time.Now(), ServiceName, installment.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense installment: %w", err)
	}
	installment.ID = existingID
	return nil
}

// GetInstallmentByExpenseAndDate retrieves an installment by expense ID and spending date (YYYY/MM).
// IMPORTANT: This method is critical for preventing duplicate installment records
// when processing invoice/savings expenses. It checks if an installment already exists
// for a specific expense and due date combination.
func (r *ExpenseInstallmentRepository) GetInstallmentByExpenseAndDate(expenseID int64, spendingDate string) (*models.ExpenseInstallment, error) {
	installment := &models.ExpenseInstallment{}

	// The due_date is stored as a DATE, so we compare using DATE_FORMAT to match YYYY/MM
	err := r.conn.db.QueryRow(`
		SELECT id, guid, expense_id, amount, paid_amount, id_status, due_date,
			created_at, created_by, updated_at, updated_by
		FROM expense_installment
		WHERE expense_id = ? AND DATE_FORMAT(due_date, '%Y/%m') = ?`,
		expenseID, spendingDate,
	).Scan(
		&installment.ID, &installment.GUID, &installment.ExpenseID, &installment.Amount, &installment.PaidAmount,
		&installment.IDStatus, &installment.DueDate,
		&installment.CreatedAt, &installment.CreatedBy, &installment.UpdatedAt, &installment.UpdatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get installment by expense/date: %w", err)
	}
	return installment, nil
}

// FinancialInstitutionRepository handles financial institution database operations
type FinancialInstitutionRepository struct {
	conn *Connection
}

// NewFinancialInstitutionRepository creates a new FinancialInstitutionRepository
func NewFinancialInstitutionRepository(conn *Connection) *FinancialInstitutionRepository {
	return &FinancialInstitutionRepository{conn: conn}
}

// UpsertFinancialInstitution inserts or updates a financial institution
func (r *FinancialInstitutionRepository) UpsertFinancialInstitution(fi *models.FinancialInstitution) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM financial_institution WHERE guid = ?", fi.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO financial_institution (guid, source_id, user_id, name, fl_credit_card, fl_money_movement, fl_investment,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			fi.GUID, fi.SourceID, fi.UserID, fi.Name, fi.FlCreditCard, fi.FlMoneyMovement, fi.FlInvestment,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert financial institution: %w", err)
		}
		id, _ := result.LastInsertId()
		fi.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check financial institution existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE financial_institution SET source_id = ?, user_id = ?, name = ?, fl_credit_card = ?, fl_money_movement = ?, fl_investment = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		fi.SourceID, fi.UserID, fi.Name, fi.FlCreditCard, fi.FlMoneyMovement, fi.FlInvestment,
		time.Now(), ServiceName, fi.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update financial institution: %w", err)
	}
	fi.ID = existingID
	return nil
}

// AdditionalBalanceRepository handles additional balance database operations
type AdditionalBalanceRepository struct {
	conn *Connection
}

// NewAdditionalBalanceRepository creates a new AdditionalBalanceRepository
func NewAdditionalBalanceRepository(conn *Connection) *AdditionalBalanceRepository {
	return &AdditionalBalanceRepository{conn: conn}
}

// UpsertAdditionalBalance inserts or updates an additional balance
func (r *AdditionalBalanceRepository) UpsertAdditionalBalance(ab *models.AdditionalBalance) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM additional_balance WHERE guid = ?", ab.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO additional_balance (guid, source_id, user_id, spending_date__YYYY_MM, amount, description,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			ab.GUID, ab.SourceID, ab.UserID, ab.SpendingDateYYYYMM, ab.Amount, ab.Description,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert additional balance: %w", err)
		}
		id, _ := result.LastInsertId()
		ab.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check additional balance existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE additional_balance SET source_id = ?, user_id = ?, spending_date__YYYY_MM = ?, amount = ?, description = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		ab.SourceID, ab.UserID, ab.SpendingDateYYYYMM, ab.Amount, ab.Description,
		time.Now(), ServiceName, ab.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update additional balance: %w", err)
	}
	ab.ID = existingID
	return nil
}

// BalanceHistoryRepository handles balance history database operations
type BalanceHistoryRepository struct {
	conn *Connection
}

// NewBalanceHistoryRepository creates a new BalanceHistoryRepository
func NewBalanceHistoryRepository(conn *Connection) *BalanceHistoryRepository {
	return &BalanceHistoryRepository{conn: conn}
}

// UpsertBalanceHistory inserts or updates a balance history record
func (r *BalanceHistoryRepository) UpsertBalanceHistory(bh *models.BalanceHistory) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM balance_history WHERE guid = ?", bh.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO balance_history (guid, source_id, user_id, spending_date__YYYY_MM, amount, last_month_amount, monthly_income,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			bh.GUID, bh.SourceID, bh.UserID, bh.SpendingDateYYYYMM, bh.Amount, bh.LastMonthAmount, bh.MonthlyIncome,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert balance history: %w", err)
		}
		id, _ := result.LastInsertId()
		bh.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check balance history existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE balance_history SET source_id = ?, user_id = ?, spending_date__YYYY_MM = ?, amount = ?, last_month_amount = ?, monthly_income = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		bh.SourceID, bh.UserID, bh.SpendingDateYYYYMM, bh.Amount, bh.LastMonthAmount, bh.MonthlyIncome,
		time.Now(), ServiceName, bh.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update balance history: %w", err)
	}
	bh.ID = existingID
	return nil
}

// ExpenseAutomaticWorkflowRepository handles expense automatic workflow database operations
type ExpenseAutomaticWorkflowRepository struct {
	conn *Connection
}

// NewExpenseAutomaticWorkflowRepository creates a new ExpenseAutomaticWorkflowRepository
func NewExpenseAutomaticWorkflowRepository(conn *Connection) *ExpenseAutomaticWorkflowRepository {
	return &ExpenseAutomaticWorkflowRepository{conn: conn}
}

// UpsertExpenseAutomaticWorkflow inserts or updates an expense automatic workflow
func (r *ExpenseAutomaticWorkflowRepository) UpsertExpenseAutomaticWorkflow(eaw *models.ExpenseAutomaticWorkflow) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM expense_automatic_workflow WHERE guid = ?", eaw.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO expense_automatic_workflow (guid, source_id, user_id, base64_image, description, extracted_expense_content_from_image,
				spending_date__YYYY_MM, sync_processed_date, id_sync_status, processing_message, created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			eaw.GUID, eaw.SourceID, eaw.UserID, eaw.Base64Image, eaw.Description, eaw.ExtractedExpenseContentFromImage,
			eaw.SpendingDateYYYYMM, eaw.SyncProcessedDate, eaw.IDSyncStatus, eaw.ProcessingMessage, time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert expense automatic workflow: %w", err)
		}
		id, _ := result.LastInsertId()
		eaw.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check expense automatic workflow existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE expense_automatic_workflow SET source_id = ?, user_id = ?, base64_image = ?, description = ?, extracted_expense_content_from_image = ?,
			spending_date__YYYY_MM = ?, sync_processed_date = ?, id_sync_status = ?, processing_message = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		eaw.SourceID, eaw.UserID, eaw.Base64Image, eaw.Description, eaw.ExtractedExpenseContentFromImage,
		eaw.SpendingDateYYYYMM, eaw.SyncProcessedDate, eaw.IDSyncStatus, eaw.ProcessingMessage,
		time.Now(), ServiceName, eaw.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense automatic workflow: %w", err)
	}
	eaw.ID = existingID
	return nil
}

// ServicePaymentRepository handles service payment database operations
type ServicePaymentRepository struct {
	conn *Connection
}

// NewServicePaymentRepository creates a new ServicePaymentRepository
func NewServicePaymentRepository(conn *Connection) *ServicePaymentRepository {
	return &ServicePaymentRepository{conn: conn}
}

// UpsertServicePayment inserts or updates a service payment
func (r *ServicePaymentRepository) UpsertServicePayment(sp *models.ServicePayment) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM service_payment WHERE guid = ?", sp.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO service_payment (guid, source_id, user_id, service_payment_date, service_payment_type_id,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			sp.GUID, sp.SourceID, sp.UserID, sp.ServicePaymentDate, sp.ServicePaymentTypeID,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert service payment: %w", err)
		}
		id, _ := result.LastInsertId()
		sp.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check service payment existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE service_payment SET source_id = ?, user_id = ?, service_payment_date = ?, service_payment_type_id = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		sp.SourceID, sp.UserID, sp.ServicePaymentDate, sp.ServicePaymentTypeID,
		time.Now(), ServiceName, sp.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update service payment: %w", err)
	}
	sp.ID = existingID
	return nil
}

// ExpenseAutomaticWorkflowPreSavedDescriptionRepository handles expense automatic workflow pre-saved description database operations
type ExpenseAutomaticWorkflowPreSavedDescriptionRepository struct {
	conn *Connection
}

// NewExpenseAutomaticWorkflowPreSavedDescriptionRepository creates a new ExpenseAutomaticWorkflowPreSavedDescriptionRepository
func NewExpenseAutomaticWorkflowPreSavedDescriptionRepository(conn *Connection) *ExpenseAutomaticWorkflowPreSavedDescriptionRepository {
	return &ExpenseAutomaticWorkflowPreSavedDescriptionRepository{conn: conn}
}

// UpsertExpenseAutomaticWorkflowPreSavedDescription inserts or updates a pre-saved description
func (r *ExpenseAutomaticWorkflowPreSavedDescriptionRepository) UpsertExpenseAutomaticWorkflowPreSavedDescription(desc *models.ExpenseAutomaticWorkflowPreSavedDescription) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM expense_automatic_workflow_pre_saved_description WHERE guid = ?", desc.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO expense_automatic_workflow_pre_saved_description (guid, source_id, user_id, description, created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?)`,
			desc.GUID, desc.SourceID, desc.UserID, desc.Description, time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert pre-saved description: %w", err)
		}
		id, _ := result.LastInsertId()
		desc.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check pre-saved description existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE expense_automatic_workflow_pre_saved_description SET source_id = ?, user_id = ?, description = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		desc.SourceID, desc.UserID, desc.Description, time.Now(), ServiceName, desc.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update pre-saved description: %w", err)
	}
	desc.ID = existingID
	return nil
}

// SystemSettingsRepository handles system settings database operations
type SystemSettingsRepository struct {
	conn *Connection
}

// NewSystemSettingsRepository creates a new SystemSettingsRepository
func NewSystemSettingsRepository(conn *Connection) *SystemSettingsRepository {
	return &SystemSettingsRepository{conn: conn}
}

// UpsertSystemSettings inserts or updates system settings
func (r *SystemSettingsRepository) UpsertSystemSettings(ss *models.SystemSettings) error {
	var existingID int64
	err := r.conn.db.QueryRow("SELECT id FROM system_settings WHERE guid = ?", ss.GUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		result, err := r.conn.db.Exec(`
			INSERT INTO system_settings (guid, source_id, fl_block_user_registration, fl_maintenance, json_sync_metadata,
				created_at, created_by)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			ss.GUID, ss.SourceID, ss.FlBlockUserRegistration, ss.FlMaintenance, ss.JSONSyncMetadata,
			time.Now(), ServiceName,
		)
		if err != nil {
			return fmt.Errorf("failed to insert system settings: %w", err)
		}
		id, _ := result.LastInsertId()
		ss.ID = id
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check system settings existence: %w", err)
	}

	_, err = r.conn.db.Exec(`
		UPDATE system_settings SET source_id = ?, fl_block_user_registration = ?, fl_maintenance = ?, json_sync_metadata = ?,
			updated_at = ?, updated_by = ?
		WHERE guid = ?`,
		ss.SourceID, ss.FlBlockUserRegistration, ss.FlMaintenance, ss.JSONSyncMetadata,
		time.Now(), ServiceName, ss.GUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update system settings: %w", err)
	}
	ss.ID = existingID
	return nil
}

// GenerateGUID generates a new UUID
func GenerateGUID() string {
	return uuid.New().String()
}
