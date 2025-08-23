// models/estimate.go contains our structs that mirror our database tables as well as our EstimateStatuses.
// the purpose of the functions in this file are to do our basic CRUD operators and any multirow results from the
// database pertaining to our estimate objects.

package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// EstimateStatus and its const values are essentially a way for us to enumerate the statuses without having to use
// string literals throughout our code when a user is progressing through the project.
// StatusDraft				-	1
// StatusAwaitingPayment	-	2
// StatusPaid				-	3
// StatusAwaitingContractor -	4
// StatusInProgress 		- 	5
// StatusCompleted			- 	6
type EstimateStatus int

const (
	StatusDraft EstimateStatus = iota + 1
	// StatusAwaitingPayment - waiting for the customer to pay for the invoice
	StatusAwaitingPayment
	// StatusPaid - customer has paid
	StatusPaid
	// StatusAwaitingContractor - Awaiting Contractor agreeing to completing the job
	StatusAwaitingContractor
	// StatusInProgress - the status of physical work being done as in literal work IN PROGRESS.
	StatusInProgress
	// StatusDone - all things complete job is done.
	StatusCompleted
)

func (s EstimateStatus) String() string {
	switch s {
	case StatusDraft:
		return "Draft"
	case StatusAwaitingPayment:
		return "Awaiting Customer Payment"
	case StatusPaid:
		return "Paid"
	case StatusAwaitingContractor:
		return "Awaiting Contractor Agreement"
	case StatusInProgress:
		return "In Progress"
	case StatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

// Estimate Struct is all the values held within the Estimate Database Object. Only values that cannot be null within
// the DB are EstimateID and CreatedBy
type Estimate struct {
	EstimateID       int // primary key
	CustomerID       int
	CreatedBy        int // surveyor who creates the estimate.
	Status           EstimateStatus
	CreatedAt        time.Time
	KitchenLengthFt  float32
	KitchenWidthFt   float32
	KitchenHeightFt  float32
	DoorWidthInches  float32
	DoorHeightInches float32
	FlooringType     string
	HasIsland        bool
	Street           string
	City             string
	State            string
	Zip              string
}

// EstimateModel wraps our sql.DB connection and allows for methods like Insert, Get, and Delete to work for estimates.
type EstimateModel struct {
	DB *sql.DB
}

// Insert creates a new estimate in the database and sets e.EstimateID.
func (m *EstimateModel) Insert(e *Estimate) error {
	stmt := `INSERT INTO estimates (
                customer_id, created_by, status, created_at,
                kitchen_length_ft, kitchen_width_ft, kitchen_height_ft,
                door_width_inches, door_height_inches, flooring_type, has_island,
                street, city, state, zip
             )
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
             RETURNING estimate_id`

	err := m.DB.QueryRow(stmt,
		e.CustomerID, e.CreatedBy, e.Status, e.CreatedAt,
		e.KitchenLengthFt, e.KitchenWidthFt, e.KitchenHeightFt,
		e.DoorWidthInches, e.DoorHeightInches, e.FlooringType, e.HasIsland,
		e.Street, e.City, e.State, e.Zip,
	).Scan(&e.EstimateID)

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an estimate by ID or returns ErrNoRecord if none found.
func (m *EstimateModel) Get(id int) (Estimate, error) {
	var estimate Estimate

	stmt := `SELECT estimate_id, customer_id, created_by, status, created_at,
       kitchen_length_ft, kitchen_width_ft, kitchen_height_ft,
       door_width_inches, door_height_inches, flooring_type, has_island,
       street, city, state, zip 
	   FROM estimates WHERE estimate_id=$1;`

	var statusInt int
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&estimate.EstimateID, &estimate.CustomerID, &estimate.CreatedBy, &statusInt, &estimate.CreatedAt, &estimate.KitchenLengthFt, &estimate.KitchenWidthFt, &estimate.KitchenHeightFt, &estimate.DoorWidthInches, &estimate.DoorHeightInches, &estimate.FlooringType, &estimate.HasIsland, &estimate.Street, &estimate.City, &estimate.State, &estimate.Zip)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Estimate{}, ErrNoRecord
		} else {
			return Estimate{}, err
		}
	}
	estimate.Status = EstimateStatus(statusInt)

	return estimate, nil

}

// GetSurveyorsEstimates retrieves up to 10 estimate by the SurveyorID (Estimate.CreatedBy) or returns ErrNoRecord if none found
// NOTE: This function has yet to be used and could potentially be changed in the future.
func (m *EstimateModel) GetSurveyorsEstimates(surveyorID int) ([]Estimate, error) {
	var estimates []Estimate
	stmt := `SELECT estimate_id, customer_id, created_by, status, created_at,
       kitchen_length_ft, kitchen_width_ft, kitchen_height_ft,
       door_width_inches, door_height_inches, flooring_type, has_island,
       street, city, state, zip 
	   FROM estimates WHERE created_by=$1 ORDER BY created_at LIMIT $2`

	rows, err := m.DB.Query(stmt, surveyorID, 10)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no estimates were found - no rows returned")
	} else if err != nil {
		return nil, err

	}

	defer rows.Close()
	for rows.Next() {
		var estimate Estimate

		err = rows.Scan(&estimate.EstimateID, &estimate.CustomerID, &estimate.CreatedBy, &estimate.CreatedAt,
			&estimate.Status, &estimate.KitchenLengthFt, &estimate.KitchenWidthFt, &estimate.KitchenHeightFt, &estimate.DoorWidthInches, &estimate.DoorHeightInches, &estimate.HasIsland, &estimate.Street, &estimate.City, &estimate.State, &estimate.Zip)
		if err != nil {
			return nil, err
		}

		estimates = append(estimates, estimate)
	}

	return estimates, nil
}

// Update updates an estimate record in the database. The update sets ALL values in the table.
// this function only returns errors.
func (m *EstimateModel) Update(e *Estimate) error {
	stmt := `UPDATE estimates
             SET customer_id = $2,status = $3,kitchen_length_ft = $4,kitchen_width_ft = $5,kitchen_height_ft = $6,door_width_inches = $7,door_height_inches = $8,flooring_type = $9,has_island = $10,street = $11,city = $12,state = $13,zip = $14
             WHERE estimate_id = $1;`

	result, err := m.DB.Exec(stmt,
		e.CustomerID, e.CreatedBy, e.Status, e.CreatedAt,
		e.KitchenLengthFt, e.KitchenWidthFt, e.KitchenHeightFt,
		e.DoorWidthInches, e.DoorHeightInches, e.FlooringType, e.HasIsland,
		e.Street, e.City, e.State, e.Zip,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected for estimate ID %d", e.EstimateID)
	}

	return nil
}

// Delete removes an Estimate from the table based off of the ID.
// This only returns Errors.
func (m *EstimateModel) Delete(id int) error {

	stmt := `DELETE FROM estimates WHERE estimate_id=$1`

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
