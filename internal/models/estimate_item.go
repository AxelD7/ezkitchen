package models

import (
	"database/sql"
	"errors"
)

type EstimateItem struct {
	LineItemID int
	EstimateID int
	ProductID  int
	Quantity   int
}

type ProductTemplateData struct {
	Product      Product
	EstimateItem EstimateItem
}

type EstimateItemModel struct {
	DB *sql.DB
}

func (m *EstimateItemModel) Insert(estimateItem *EstimateItem) error {

	stmt := `INSERT INTO estimate_items (estimate_id, product_id, quantity)
	VALUES ($1,$2,$3) RETURNING line_item_id`

	err := m.DB.QueryRow(stmt, estimateItem.EstimateID, estimateItem.ProductID, estimateItem.Quantity).Scan(&estimateItem.LineItemID)
	if err != nil {
		return err
	}
	return nil

}

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

func (m *EstimateItemModel) GetByEstimateID(id int) ([]ProductTemplateData, error) {
	var estimateProducts []ProductTemplateData
	stmt := `SELECT ei.line_item_id, ei.product_id, ei.quantity, p.name, p.description, p.category, p.subcategory, p.color,  p.unit_price FROM estimate_items ei INNER JOIN products p on ei.product_id = p.product_id WHERE estimate_id=$1`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var estimateProduct ProductTemplateData

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
