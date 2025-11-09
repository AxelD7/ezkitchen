package integration_test

import (
	"errors"
	"ezkitchen/internal/models"
	"testing"
)

func TestEstimateItemInsertAndGet(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	estimate := createTestEstimate(t, customer.ID, surveyor.ID)
	product := createTestProduct(t, surveyor.ID)

	item := &models.EstimateItem{
		EstimateID: estimate.EstimateID,
		ProductID:  product.ProductID,
		Quantity:   3,
	}

	if err := estimateItemModel.Insert(item); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if item.LineItemID == 0 {
		t.Fatalf("Expected non-zero LineItemID after insert")
	}

	got, err := estimateItemModel.GetByLineItemID(item.LineItemID)
	if err != nil {
		t.Fatalf("GetByLineItemID failed: %v", err)
	}
	if got.Quantity != item.Quantity {
		t.Errorf("Quantity mismatch: got %d want %d", got.Quantity, item.Quantity)
	}
	if got.EstimateID != estimate.EstimateID {
		t.Errorf("EstimateID mismatch: got %d want %d", got.EstimateID, estimate.EstimateID)
	}
	if got.ProductID != product.ProductID {
		t.Errorf("ProductID mismatch: got %d want %d", got.ProductID, product.ProductID)
	}
}

func TestEstimateItemUpdate(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	estimate := createTestEstimate(t, customer.ID, surveyor.ID)
	product := createTestProduct(t, surveyor.ID)

	item := &models.EstimateItem{
		EstimateID: estimate.EstimateID,
		ProductID:  product.ProductID,
		Quantity:   2,
	}
	if err := estimateItemModel.Insert(item); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	item.Quantity = 5
	if err := estimateItemModel.Update(*item); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := estimateItemModel.GetByLineItemID(item.LineItemID)
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}
	if got.Quantity != item.Quantity {
		t.Errorf("Expected Quantity %d, got %d", item.Quantity, got.Quantity)
	}
}

func TestEstimateItemDelete(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	estimate := createTestEstimate(t, customer.ID, surveyor.ID)
	product := createTestProduct(t, surveyor.ID)

	item := &models.EstimateItem{
		EstimateID: estimate.EstimateID,
		ProductID:  product.ProductID,
		Quantity:   4,
	}
	if err := estimateItemModel.Insert(item); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if err := estimateItemModel.Delete(item.LineItemID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Deleting again should return ErrNoRecord
	err := estimateItemModel.Delete(item.LineItemID)
	if err == nil {
		t.Fatal("Expected ErrNoRecord when deleting already deleted item, got nil")
	}
	if !errors.Is(err, models.ErrNoRecord) {
		t.Fatalf("Expected ErrNoRecord, got %v", err)
	}
}

func TestEstimateItemGetByEstimateID(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	customer := createTestUser(t, "John Smith", "john@example.com", "customer")
	surveyor := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	estimate := createTestEstimate(t, customer.ID, surveyor.ID)
	product1 := createTestProduct(t, surveyor.ID)
	product2 := &models.Product{
		Name:        "Marble Countertop",
		Description: "Elegant marble design",
		Category:    "Countertops",
		Subcategory: "Marble",
		Color:       "White",
		UnitPrice:   30000,
		CreatedBy:   surveyor.ID,
	}
	if err := productModel.Insert(product2); err != nil {
		t.Fatalf("insert product2 failed: %v", err)
	}

	item1 := &models.EstimateItem{EstimateID: estimate.EstimateID, ProductID: product1.ProductID, Quantity: 2}
	item2 := &models.EstimateItem{EstimateID: estimate.EstimateID, ProductID: product2.ProductID, Quantity: 1}
	if err := estimateItemModel.Insert(item1); err != nil {
		t.Fatalf("Insert item1 failed: %v", err)
	}
	if err := estimateItemModel.Insert(item2); err != nil {
		t.Fatalf("Insert item2 failed: %v", err)
	}

	items, err := estimateItemModel.GetByEstimateID(estimate.EstimateID)
	if err != nil {
		t.Fatalf("GetByEstimateID failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 estimate items, got %d", len(items))
	}
}
