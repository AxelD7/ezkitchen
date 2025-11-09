package integration_test

import (
	"errors"
	"ezkitchen/internal/models"
	"testing"
	"time"
)

type testUser struct{ ID int }

func createTestUser(t *testing.T, name, email, role string) testUser {
	t.Helper()
	stmt := `INSERT INTO users (name, email, password_hash, role, created_at)
	         VALUES ($1, $2, 'testhash', $3, NOW())
	         RETURNING user_id`
	var id int
	if err := testDB.QueryRow(stmt, name, email, role).Scan(&id); err != nil {
		t.Fatalf("insert user failed: %v", err)
	}
	return testUser{ID: id}
}

func createTestEstimate(t *testing.T, customerID, surveyorID int) *models.Estimate {
	t.Helper()
	e := &models.Estimate{
		CustomerID:        customerID,
		CreatedBy:         surveyorID,
		Status:            models.StatusDraft,
		CreatedAt:         time.Now(),
		KitchenLengthInch: 120,
		KitchenWidthInch:  100,
		KitchenHeightInch: 96,
		DoorWidthInch:     36,
		DoorHeightInch:    80,
		Street:            "123 Test Ave",
		City:              "Detroit",
		State:             "MI",
		Zip:               "48201",
	}

	if err := estimateModel.Insert(e); err != nil {
		t.Fatalf("insert estimate failed: %v", err)
	}
	return e
}

func TestEstimateInsertAndGet(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "babyboss@example.com", "surveyor")

	e := createTestEstimate(t, customer.ID, surveyor.ID)

	if e.EstimateID == 0 {
		t.Fatalf("Expected a non-zero EstimateID after insert")
	}

	got, err := estimateModel.Get(e.EstimateID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got.CustomerID != e.CustomerID {
		t.Errorf("CustomerID mismatch: got %d want %d", got.CustomerID, e.CustomerID)
	}
	if got.CreatedBy != e.CreatedBy {
		t.Errorf("CreatedBy mismatch: got %d want %d", got.CreatedBy, e.CreatedBy)
	}
	if got.Status != e.Status {
		t.Errorf("Status mismatch: got %v want %v", got.Status, e.Status)
	}
	if got.City != e.City {
		t.Errorf("City mismatch: got %q want %q", got.City, e.City)
	}
}

func TestEstimateUpdate(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "babyboss@example.com", "surveyor")

	e := createTestEstimate(t, customer.ID, surveyor.ID)

	e.Status = models.StatusInProgress
	e.DoorWidthInch = 46
	e.City = "Royal Oak"
	e.Zip = "48120"

	if err := estimateModel.Update(e); err != nil {
		t.Fatalf("Update Failed: %v", err)
	}

	got, err := estimateModel.Get(e.EstimateID)
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}

	if got.Status != e.Status {
		t.Errorf("Expected Status %v got %v", e.Status, got.Status)
	}
	if got.KitchenWidthInch != e.KitchenWidthInch {
		t.Errorf("Expected KitchenWidthInch %.1f got %.1f", e.KitchenWidthInch, got.KitchenWidthInch)
	}
	if got.City != e.City {
		t.Errorf("Expected City %q got %q", e.City, got.City)
	}
	if got.Zip != e.Zip {
		t.Errorf("Expected Zip %q got %q", e.Zip, got.Zip)
	}
}

func TestEstimateDelete(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "babyboss@example.com", "surveyor")

	e := createTestEstimate(t, customer.ID, surveyor.ID)

	if err := estimateModel.Delete(e.EstimateID); err != nil {
		t.Fatalf("Delete Failed: %v", err)
	}

	// Deleting again should return ErrNoRecord
	err := estimateModel.Delete(e.EstimateID)
	if err == nil {
		t.Fatal("Expected ErrNoRecord when deleting an already deleted record, got nil")
	}
	if !errors.Is(err, models.ErrNoRecord) {
		t.Fatalf("Expected ErrNoRecord, got %v", err)
	}
}
