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

func (m *EstimateItemModel) GetByEstimateID(id int) ([]EstimateItem, error) {
	var estimateItems []EstimateItem
	stmt := `SELECT line_item_id, estimate_id, product_id, quantity FROM estimate_items WHERE estimate_id=$1`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var estimateItem EstimateItem

		err := rows.Scan(&estimateItem.LineItemID, &estimateItem.EstimateID, &estimateItem.ProductID, &estimateItem.Quantity)
		if err != nil {
			return nil, err
		}

		estimateItems = append(estimateItems, estimateItem)

	}

	return estimateItems, nil
}

func (m *EstimateItemModel) Update(estimateItem EstimateItem) error {
	stmt := `UPDATE estimate_items SET estimate_id=$2, product_id=$3, quantity=$4 WHERE line_item_id=$1`
	result, err := m.DB.Exec(stmt, estimateItem.LineItemID, estimateItem.EstimateID, estimateItem.ProductID, estimateItem.Quantity)
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
