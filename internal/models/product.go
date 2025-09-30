package models

import (
	"database/sql"
	"errors"
)

type Product struct {
	ProductID   int
	Name        string
	Description string
	Category    string
	UnitPrice   float64
	CreatedBy   int
}

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) Insert(p *Product) error {
	stmt := `INSERT INTO products (name, description, category, unit_price, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id`

	err := m.DB.QueryRow(stmt, p.Name, p.Description, p.Category, p.UnitPrice, p.CreatedBy).Scan(&p.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func (m *ProductModel) Get(id int) (Product, error) {
	stmt := `SELECT product_id, name, description, category, unit_price, created_by
		FROM products WHERE product_id=$1`

	var p Product
	err := m.DB.QueryRow(stmt, id).Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.UnitPrice, &p.CreatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, ErrNoRecord
		}
		return Product{}, err
	}

	return p, nil
}

func (m *ProductModel) GetByProductCategory(category string) ([]Product, error) {
	stmt := `SELECT product_id, name, description, category, unit_price, created_by FROM products WHERE category=$1`
	var products []Product

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductID, &product.Name, &product.Description, &product.Category, &product.UnitPrice, &product.CreatedBy)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (m *ProductModel) Update(p *Product) error {
	stmt := `UPDATE products
		SET name=$2, description=$3, category=$4, unit_price=$5, created_by=$6
		WHERE product_id=$1`

	result, err := m.DB.Exec(stmt, p.ProductID, p.Name, p.Description, p.Category, p.UnitPrice, p.CreatedBy)
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

func (m *ProductModel) Delete(id int) error {
	stmt := `DELETE FROM products WHERE product_id=$1`

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
