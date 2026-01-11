package mariadb

import (
	"testing"
)

func TestGenerateGUID(t *testing.T) {
	guid1 := GenerateGUID()
	guid2 := GenerateGUID()

	// Check that GUIDs are not empty
	if guid1 == "" {
		t.Error("GenerateGUID() returned empty string")
	}

	// Check that GUIDs are unique
	if guid1 == guid2 {
		t.Error("GenerateGUID() returned same GUID twice")
	}

	// Check that GUID has correct format (UUID v4)
	if len(guid1) != 36 {
		t.Errorf("GenerateGUID() returned GUID with length %d, want 36", len(guid1))
	}

	// Check for dashes at correct positions
	if guid1[8] != '-' || guid1[13] != '-' || guid1[18] != '-' || guid1[23] != '-' {
		t.Errorf("GenerateGUID() returned GUID with incorrect format: %s", guid1)
	}
}

func TestNewUserRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewUserRepository(conn)

	if repo == nil {
		t.Error("NewUserRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewUserRepository() didn't set connection correctly")
	}
}

func TestNewExpenseRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewExpenseRepository(conn)

	if repo == nil {
		t.Error("NewExpenseRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewExpenseRepository() didn't set connection correctly")
	}
}

func TestNewExpenseInstallmentRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewExpenseInstallmentRepository(conn)

	if repo == nil {
		t.Error("NewExpenseInstallmentRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewExpenseInstallmentRepository() didn't set connection correctly")
	}
}

func TestNewFinancialInstitutionRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewFinancialInstitutionRepository(conn)

	if repo == nil {
		t.Error("NewFinancialInstitutionRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewFinancialInstitutionRepository() didn't set connection correctly")
	}
}

func TestNewAdditionalBalanceRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewAdditionalBalanceRepository(conn)

	if repo == nil {
		t.Error("NewAdditionalBalanceRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewAdditionalBalanceRepository() didn't set connection correctly")
	}
}

func TestNewBalanceHistoryRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewBalanceHistoryRepository(conn)

	if repo == nil {
		t.Error("NewBalanceHistoryRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewBalanceHistoryRepository() didn't set connection correctly")
	}
}

func TestNewExpenseAutomaticWorkflowRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewExpenseAutomaticWorkflowRepository(conn)

	if repo == nil {
		t.Error("NewExpenseAutomaticWorkflowRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewExpenseAutomaticWorkflowRepository() didn't set connection correctly")
	}
}

func TestNewServicePaymentRepository(t *testing.T) {
	conn := &Connection{db: nil}
	repo := NewServicePaymentRepository(conn)

	if repo == nil {
		t.Error("NewServicePaymentRepository() returned nil")
	}

	if repo.conn != conn {
		t.Error("NewServicePaymentRepository() didn't set connection correctly")
	}
}
