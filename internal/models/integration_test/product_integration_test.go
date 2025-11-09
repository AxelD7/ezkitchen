package integration_test

import (
	"errors"
	"ezkitchen/internal/models"
	"testing"
)

func createTestProduct(t *testing.T, createdBy int) *models.Product {
	t.Helper()
	p := &models.Product{
		Name:        "Test Product",
		Description: "A durable countertop",
		Category:    "Countertops",
		Subcategory: "Granite",
		Color:       "Black",
		UnitPrice:   25000,
		Length:      72.0,
		Width:       24.0,
		Height:      1.5,
		CreatedBy:   createdBy,
	}
	if err := productModel.Insert(p); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	return p
}

func TestProductInsertAndGet(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	user := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	p := createTestProduct(t, user.ID)
	if p.ProductID == 0 {
		t.Fatalf("expected non-zero ProductID after insert")
	}

	got, err := productModel.Get(p.ProductID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got.Name != p.Name {
		t.Errorf("Name mismatch: got %q want %q", got.Name, p.Name)
	}
	if got.Color != p.Color {
		t.Errorf("Color mismatch: got %q want %q", got.Color, p.Color)
	}
	if got.UnitPrice != p.UnitPrice {
		t.Errorf("UnitPrice mismatch: got %d want %d", got.UnitPrice, p.UnitPrice)
	}
	if got.CreatedBy != p.CreatedBy {
		t.Errorf("CreatedBy mismatch: got %d want %d", got.CreatedBy, p.CreatedBy)
	}
}

func TestProductUpdate(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	user := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")
	p := createTestProduct(t, user.ID)

	p.Name = "Updated Product"
	p.Color = "White"
	p.UnitPrice = 27500
	p.Length = 60.0

	if err := productModel.Update(p); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := productModel.Get(p.ProductID)
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}

	if got.Name != p.Name {
		t.Errorf("Expected Name %q, got %q", p.Name, got.Name)
	}
	if got.Color != p.Color {
		t.Errorf("Expected Color %q, got %q", p.Color, got.Color)
	}
	if got.UnitPrice != p.UnitPrice {
		t.Errorf("Expected UnitPrice %d, got %d", p.UnitPrice, got.UnitPrice)
	}
	if got.Length != p.Length {
		t.Errorf("Expected Length %.1f, got %.1f", p.Length, got.Length)
	}
}

func TestProductDelete(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	user := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")
	p := createTestProduct(t, user.ID)

	if err := productModel.Delete(p.ProductID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Second delete should return ErrNoRecord
	err := productModel.Delete(p.ProductID)
	if err == nil {
		t.Fatal("Expected ErrNoRecord when deleting already deleted product, got nil")
	}
	if !errors.Is(err, models.ErrNoRecord) {
		t.Fatalf("Expected ErrNoRecord, got %v", err)
	}
}

func TestProductFilter(t *testing.T) {
	t.Cleanup(func() { resetDB(t) })

	user := createTestUser(t, "Daniel Surveyor", "boss@example.com", "surveyor")

	p1 := &models.Product{
		Name:        "Granite Black",
		Description: "Luxury granite top",
		Category:    "Countertops",
		Subcategory: "Granite",
		Color:       "Black",
		UnitPrice:   20000,
		CreatedBy:   user.ID,
	}
	if err := productModel.Insert(p1); err != nil {
		t.Fatalf("Insert product 1 failed: %v", err)
	}

	p2 := &models.Product{
		Name:        "Marble White",
		Description: "Elegant marble top",
		Category:    "Countertops",
		Subcategory: "Marble",
		Color:       "White",
		UnitPrice:   30000,
		CreatedBy:   user.ID,
	}
	if err := productModel.Insert(p2); err != nil {
		t.Fatalf("Insert product 2 failed: %v", err)
	}

	results, err := productModel.GetByProductFilter("Countertops", "", "")
	if err != nil {
		t.Fatalf("GetByProductFilter failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results for category 'Countertops', got %d", len(results))
	}

	results, err = productModel.GetByProductFilter("", "Marble", "")
	if err != nil {
		t.Fatalf("GetByProductFilter by subcategory failed: %v", err)
	}
	if len(results) != 1 || results[0].Subcategory != "Marble" {
		t.Errorf("Expected 1 Marble product, got %+v", results)
	}

	results, err = productModel.GetByProductFilter("", "", "Black")
	if err != nil {
		t.Fatalf("GetByProductFilter by color failed: %v", err)
	}
	if len(results) != 1 || results[0].Color != "Black" {
		t.Errorf("Expected 1 Black product, got %+v", results)
	}
}
