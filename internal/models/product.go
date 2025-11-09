package models

import (
	"database/sql"
	"errors"
)

// Product represents a product record in the database.
// It contains all descriptive and dimensional data for a given item.
type Product struct {
	ProductID   int
	Name        string
	Description string
	Category    string
	Subcategory string
	Color       string
	UnitPrice   int
	Length      float32
	Width       float32
	Height      float32
	CreatedBy   int
}

// ProductModel wraps a sql.DB connection and provides methods for CRUD operations on products.
type ProductModel struct {
	DB *sql.DB
}

// Insert adds a new Product to the database and assigns the generated ProductID to the struct.
// Returns an error if the insert or Scan operation fails.
func (m *ProductModel) Insert(p *Product) error {
	stmt := `
		INSERT INTO products
			(name, description, category, subcategory, color, unit_price, length, width, height, created_by)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING product_id
	`
	return m.DB.QueryRow(stmt,
		p.Name,
		p.Description,
		p.Category,
		p.Subcategory,
		p.Color,
		p.UnitPrice,
		p.Length,
		p.Width,
		p.Height,
		p.CreatedBy,
	).Scan(&p.ProductID)
}

// Get retrieves a Product by its ProductID.
// Returns ErrNoRecord if the product does not exist.
func (m *ProductModel) Get(id int) (Product, error) {
	stmt := `
		SELECT product_id, name, description, category, subcategory, color,
		       unit_price, length, width, height, created_by
		FROM products
		WHERE product_id=$1
	`
	var p Product
	err := m.DB.QueryRow(stmt, id).Scan(
		&p.ProductID,
		&p.Name,
		&p.Description,
		&p.Category,
		&p.Subcategory,
		&p.Color,
		&p.UnitPrice,
		&p.Length,
		&p.Width,
		&p.Height,
		&p.CreatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Product{}, ErrNoRecord
		}
		return Product{}, err
	}
	return p, nil
}

// GetByProductFilter retrieves a list of Products filtered by category, subcategory, or color.
// Empty string values are treated as wildcards (i.e., no filter applied).
// Returns a slice of Product or an error.
func (m *ProductModel) GetByProductFilter(category, subcategory, color string) ([]Product, error) {
	stmt := `
		SELECT product_id, name, description, category, subcategory, color,
		       unit_price, length, width, height, created_by
		FROM products
		WHERE ($1='' OR category=$1)
		AND ($2='' OR subcategory=$2)
		AND ($3='' OR color=$3)
	`

	rows, err := m.DB.Query(stmt, category, subcategory, color)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ProductID,
			&p.Name,
			&p.Description,
			&p.Category,
			&p.Subcategory,
			&p.Color,
			&p.UnitPrice,
			&p.Length,
			&p.Width,
			&p.Height,
			&p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// Update modifies an existing Product record in the database.
// All fields are updated based on the provided Product struct.
// Returns ErrNoRecord if the record does not exist.
func (m *ProductModel) Update(p *Product) error {
	stmt := `
		UPDATE products
		SET name=$2, description=$3, category=$4, subcategory=$5, color=$6,
		    unit_price=$7, length=$8, width=$9, height=$10, created_by=$11
		WHERE product_id=$1
	`
	result, err := m.DB.Exec(stmt,
		p.ProductID,
		p.Name,
		p.Description,
		p.Category,
		p.Subcategory,
		p.Color,
		p.UnitPrice,
		p.Length,
		p.Width,
		p.Height,
		p.CreatedBy,
	)
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

// Delete removes a Product from the database by its ProductID.
// Returns ErrNoRecord if the specified record does not exist.
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
