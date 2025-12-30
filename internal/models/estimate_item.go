package models

import (
	"database/sql"
	"errors"
)

// EstimateItem represents a single product entry within an estimate.
type EstimateItem struct {
	LineItemID int
	EstimateID int
	ProductID  int
	Quantity   int
}

// EstimateProduct combines an EstimateItem with its associated Product data.
type EstimateProduct struct {
	Product      Product
	EstimateItem EstimateItem
}

// EstimateItemModel wraps database operations for estimate_items.
type EstimateItemModel struct {
	DB *sql.DB
}

// Insert adds a new EstimateItem to the database.
// The provided EstimateItem must have valid EstimateID and ProductID fields.
// Returns an error if the insert operation or Scan fails.
func (m *EstimateItemModel) Insert(estimateItem *EstimateItem) error {

	stmt := `INSERT INTO estimate_items (estimate_id, product_id, quantity)
	VALUES ($1,$2,$3) RETURNING line_item_id`

	err := m.DB.QueryRow(stmt, estimateItem.EstimateID, estimateItem.ProductID, estimateItem.Quantity).Scan(&estimateItem.LineItemID)
	if err != nil {
		return err
	}
	return nil

}

// GetByLineItemID retrieves an EstimateItem by its LineItemID.
// Returns the EstimateItem object, or ErrNoRecord if no matching record exists.
func (m *EstimateItemModel) GetByLineItemID(id int) (EstimateItem, error) {
	var estimateItem EstimateItem
	stmt := `SELECT line_item_id, estimate_id, product_id, quantity FROM estimate_items WHERE line_item_id=$1`
	err := m.DB.QueryRow(stmt, id).Scan(&estimateItem.LineItemID, &estimateItem.EstimateID, &estimateItem.ProductID, &estimateItem.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EstimateItem{}, ErrNoRecord
		}
		return EstimateItem{}, err
	}
	return estimateItem, nil
}

func (m *EstimateItemModel) GetEstimateIDByLineItemID(lineItemID int) (int, error) {
	var estimateID int
	stmt := `SELECT estimate_id FROM estimate_items WHERE line_item_id=$1`

	err := m.DB.QueryRow(stmt, lineItemID).Scan(&estimateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		}
		return 0, err
	}

	return estimateID, nil

}

// GetByEstimateID returns all EstimateItems that belong to the specified EstimateID.
// Each returned record includes associated Product information.
// Returns a slice of EstimateProduct or an error.
func (m *EstimateItemModel) GetByEstimateID(estimateID int) ([]EstimateProduct, error) {
	var estimateProducts []EstimateProduct
	stmt := `SELECT ei.line_item_id, ei.product_id, ei.quantity, p.name, p.description, p.category, p.subcategory, p.color,  p.unit_price FROM estimate_items ei INNER JOIN products p on ei.product_id = p.product_id WHERE estimate_id=$1`
	rows, err := m.DB.Query(stmt, estimateID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var estimateProduct EstimateProduct

		err := rows.Scan(&estimateProduct.EstimateItem.LineItemID, &estimateProduct.EstimateItem.ProductID, &estimateProduct.EstimateItem.Quantity, &estimateProduct.Product.Name, &estimateProduct.Product.Description, &estimateProduct.Product.Category, &estimateProduct.Product.Subcategory, &estimateProduct.Product.Color, &estimateProduct.Product.UnitPrice)
		if err != nil {
			return nil, err
		}

		estimateProducts = append(estimateProducts, estimateProduct)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return estimateProducts, nil
}

// Update modifies the quantity of an existing EstimateItem.
// The provided EstimateItem must include a valid LineItemID.
// Returns ErrNoRecord if no record was updated.
func (m *EstimateItemModel) Update(estimateItem EstimateItem) error {
	stmt := `UPDATE estimate_items SET quantity=$2 WHERE line_item_id=$1`
	result, err := m.DB.Exec(stmt, estimateItem.LineItemID, estimateItem.Quantity)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRecord
	}

	return nil
}

// Delete removes an EstimateItem by its LineItemID.
// Returns ErrNoRecord if the record does not exist.
func (m *EstimateItemModel) Delete(id int) error {
	stmt := `DELETE FROM estimate_items WHERE line_item_id=$1`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRecord
	}

	return nil
}
