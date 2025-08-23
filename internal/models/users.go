// models/users.go defines the User domain model and its database interactions.
// This includes CRUD operations and role definitions for administrators,
// surveyors, and customers.

package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Role represents the type of user in the system.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSurveyor Role = "surveyor"
	RoleCustomer Role = "customer"
)

// User represents a customer, surveyor, or admin account in the system.
type User struct {
	UserID       int
	Name         string
	Email        string
	PasswordHash string // IMPORTANT: this is currently dummy data and awaiting user auth for full implementation.
	Phone        string
	Role         Role
	CreatedAt    time.Time
}

// UserModel wraps a sql.DB connection and provides methods for managing user tables.
type UserModel struct {
	DB *sql.DB
}

// Insert creates a new user record in the database and sets u.UserID.
// Returns an error if the insert fails.
func (m *UserModel) Insert(u *User) error {
	stmt := `INSERT INTO users (name, email, password_hash, role, phone, created_at)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING user_id`

	err := m.DB.QueryRow(stmt,
		u.Name, u.Email, u.PasswordHash, u.Role, u.Phone, u.CreatedAt,
	).Scan(&u.UserID)
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves a user by ID. Returns sql.ErrNoRows if no matching record exists.
func (m *UserModel) Get(id int) (User, error) {
	var user User
	user.UserID = id

	stmt := `SELECT name, email, password_hash, role, phone, created_at
             FROM users WHERE user_id=$1`

	row := m.DB.QueryRow(stmt, user.UserID)
	err := row.Scan(&user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, sql.ErrNoRows
	} else if err != nil {
		return User{}, err
	}
	return user, nil
}

// Update changes user fields for the given UserID.
// Returns an error if the update fails or affects no rows.
// NOTE: THIS DOES NOT CHANGE THE PASSWORD FIELD!!!!!
func (m *UserModel) Update(u *User) error {
	stmt := `UPDATE users 
             SET name=$2, email=$3, role=$4, phone=$5
             WHERE user_id=$1`

	result, err := m.DB.Exec(stmt, u.UserID, u.Name, u.Email, u.Role, u.Phone)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected for user id: %d", u.UserID)
	}

	return nil
}

// Delete removes a user by ID. Returns an error if no rows were affected.
func (m *UserModel) Delete(id int) error {
	stmt := `DELETE FROM users WHERE user_id=$1`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}
