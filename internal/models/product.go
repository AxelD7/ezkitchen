package models

import (
	"database/sql"
	"errors"
)

type Product struct {
	ProductID   int
	Name        string
	Description string
	UnitPrice   float64
	CreatedBy   int
}

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) Insert(p *Product) error {
	stmt := `INSERT INTO products (name, description, unit_price, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING product_id`

	err := m.DB.QueryRow(stmt, p.Name, p.Description, p.UnitPrice, p.CreatedBy).Scan(&p.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ProductModel) Get(id int) (Product, error) {
	stmt := `SELECT product_id, name, description, unit_price, created_by
		FROM products WHERE product_id=$1`

	var p Product
	err := m.DB.QueryRow(stmt, id).Scan(&p.ProductID, &p.Name, &p.Description, &p.UnitPrice, &p.CreatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, ErrNoRecord
		}
		return Product{}, err
	}

	return p, nil
}

func (m *ProductModel) Update(p *Product) error {
	stmt := `UPDATE products
		SET name=$2, description=$3, unit_price=$4, created_by=$5
		WHERE product_id=$1 RETURNING product_id`

	err := m.DB.QueryRow(stmt, p.ProductID, p.Name, p.Description, p.UnitPrice, p.CreatedBy).Scan(&p.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
		return err
	}

	return nil
}

func (m *ProductModel) Delete(id int) error {
	stmt := `DELETE FROM products WHERE product_id=$1 RETURNING product_id`

	err := m.DB.QueryRow(stmt, id).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}

		return err
	}

	return nil
}
