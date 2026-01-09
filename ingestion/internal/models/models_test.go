package models

import (
	"testing"
)

func TestGetDomainSeeds(t *testing.T) {
	seeds := GetDomainSeeds()

	if len(seeds) != 5 {
		t.Errorf("GetDomainSeeds() returned %d seeds, want 5", len(seeds))
	}

	// Test expense_automatic_workflow / id_sync_status
	found := false
	for _, seed := range seeds {
		if seed.Source == "expense_automatic_workflow" && seed.Type == "id_sync_status" {
			found = true
			expectedNames := []string{"pending", "success", "error"}
			if len(seed.Names) != len(expectedNames) {
				t.Errorf("id_sync_status has %d names, want %d", len(seed.Names), len(expectedNames))
			}
			for i, name := range expectedNames {
				if seed.Names[i] != name {
					t.Errorf("id_sync_status name[%d] = %s, want %s", i, seed.Names[i], name)
				}
			}
		}
	}
	if !found {
		t.Error("expense_automatic_workflow / id_sync_status not found")
	}

	// Test expense / id_status
	found = false
	for _, seed := range seeds {
		if seed.Source == "expense" && seed.Type == "id_status" {
			found = true
			expectedNames := []string{"pending", "partially_paid", "paid"}
			if len(seed.Names) != len(expectedNames) {
				t.Errorf("expense id_status has %d names, want %d", len(seed.Names), len(expectedNames))
			}
		}
	}
	if !found {
		t.Error("expense / id_status not found")
	}

	// Test expense / id_type
	found = false
	for _, seed := range seeds {
		if seed.Source == "expense" && seed.Type == "id_type" {
			found = true
			expectedNames := []string{"expense", "invoice", "savings"}
			if len(seed.Names) != len(expectedNames) {
				t.Errorf("expense id_type has %d names, want %d", len(seed.Names), len(expectedNames))
			}
		}
	}
	if !found {
		t.Error("expense / id_type not found")
	}

	// Test expense_installment / id_status
	found = false
	for _, seed := range seeds {
		if seed.Source == "expense_installment" && seed.Type == "id_status" {
			found = true
			expectedNames := []string{"pending", "partially_paid", "paid"}
			if len(seed.Names) != len(expectedNames) {
				t.Errorf("expense_installment id_status has %d names, want %d", len(seed.Names), len(expectedNames))
			}
		}
	}
	if !found {
		t.Error("expense_installment / id_status not found")
	}

	// Test service_payment / service_payment_type_id
	found = false
	for _, seed := range seeds {
		if seed.Source == "service_payment" && seed.Type == "service_payment_type_id" {
			found = true
			if len(seed.Names) != 1 || seed.Names[0] != "PayPal" {
				t.Error("service_payment_type_id should have only 'PayPal'")
			}
		}
	}
	if !found {
		t.Error("service_payment / service_payment_type_id not found")
	}
}

func TestDomainStructure(t *testing.T) {
	domain := Domain{
		ID:     1,
		GUID:   "test-guid",
		Name:   "pending",
		Type:   "id_status",
		Source: "expense",
	}

	if domain.ID != 1 {
		t.Errorf("Domain.ID = %d, want 1", domain.ID)
	}
	if domain.GUID != "test-guid" {
		t.Errorf("Domain.GUID = %s, want test-guid", domain.GUID)
	}
	if domain.Name != "pending" {
		t.Errorf("Domain.Name = %s, want pending", domain.Name)
	}
	if domain.Type != "id_status" {
		t.Errorf("Domain.Type = %s, want id_status", domain.Type)
	}
	if domain.Source != "expense" {
		t.Errorf("Domain.Source = %s, want expense", domain.Source)
	}
}

func TestUserStructure(t *testing.T) {
	user := User{
		ID:        1,
		GUID:      "user-guid",
		FirstName: "John",
		Email:     "john@example.com",
		FlAdmin:   true,
	}

	if user.ID != 1 {
		t.Errorf("User.ID = %d, want 1", user.ID)
	}
	if user.FirstName != "John" {
		t.Errorf("User.FirstName = %s, want John", user.FirstName)
	}
	if !user.FlAdmin {
		t.Error("User.FlAdmin = false, want true")
	}
}

func TestExpenseStructure(t *testing.T) {
	expense := Expense{
		ID:                 1,
		GUID:               "expense-guid",
		UserID:             1,
		SpendingDateYYYYMM: "2024-01",
		Name:               "Test Expense",
		TotalAmount:        100.50,
		TotalPaidAmount:    50.25,
	}

	if expense.SpendingDateYYYYMM != "2024-01" {
		t.Errorf("Expense.SpendingDateYYYYMM = %s, want 2024-01", expense.SpendingDateYYYYMM)
	}
	if expense.TotalAmount != 100.50 {
		t.Errorf("Expense.TotalAmount = %f, want 100.50", expense.TotalAmount)
	}
}
