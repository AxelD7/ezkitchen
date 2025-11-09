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
	// StatusCompleted - all things complete job is done.
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
	EstimateID        int // primary key
	CustomerID        int
	CreatedBy         int // surveyor who creates the estimate.
	Status            EstimateStatus
	CreatedAt         time.Time
	KitchenLengthInch float32
	KitchenWidthInch  float32
	KitchenHeightInch float32
	DoorWidthInch     float32
	DoorHeightInch    float32
	Street            string
	City              string
	State             string
	Zip               string
}

type EstimateTotals struct {
	Subtotal      int
	LaborTotal    int
	SalesTax      int
	EstimateTotal int
}

// EstimateModel wraps our sql.DB connection and allows for methods like Insert, Get, and Delete to work for estimates.
type EstimateModel struct {
	DB    *sql.DB
	Items *EstimateItemModel
}

// Insert creates a new estimate in the database and assigns the generated EstimateID to the provided struct.
// Returns an error if the insert operation or Scan fails.
func (m *EstimateModel) Insert(e *Estimate) error {
	stmt := `INSERT INTO estimates 
	(customer_id, created_by, status, created_at,
    kitchen_length_inch, kitchen_width_inch, kitchen_height_inch,
    door_width_inch, door_height_inch,
    street, city, state, zip)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING estimate_id`

	err := m.DB.QueryRow(stmt,
		e.CustomerID, e.CreatedBy, e.Status, e.CreatedAt,
		e.KitchenLengthInch, e.KitchenWidthInch, e.KitchenHeightInch,
		e.DoorWidthInch, e.DoorHeightInch, e.Street, e.City, e.State, e.Zip,
	).Scan(&e.EstimateID)

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an Estimate by its ID.
// Returns ErrNoRecord if the specified record does not exist.
func (m *EstimateModel) Get(id int) (Estimate, error) {
	var estimate Estimate

	stmt := `SELECT estimate_id, customer_id, created_by, status, created_at,
       	kitchen_length_inch, kitchen_width_inch, kitchen_height_inch,
    	door_width_inch, door_height_inch, street, city, state, zip 
	   	FROM estimates WHERE estimate_id=$1;`

	var statusInt int
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&estimate.EstimateID, &estimate.CustomerID, &estimate.CreatedBy, &statusInt, &estimate.CreatedAt, &estimate.KitchenLengthInch, &estimate.KitchenWidthInch, &estimate.KitchenHeightInch, &estimate.DoorWidthInch, &estimate.DoorHeightInch, &estimate.Street, &estimate.City, &estimate.State, &estimate.Zip)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Estimate{}, ErrNoRecord
		}
		return Estimate{}, err
	}
	estimate.Status = EstimateStatus(statusInt)

	return estimate, nil
}

// GetSurveyorsEstimates retrieves up to 10 estimates created by a specific surveyor (CreatedBy field).
// Returns ErrNoRecord or an empty slice if no estimates are found.
func (m *EstimateModel) GetSurveyorsEstimates(surveyorID int) ([]Estimate, error) {
	var estimates []Estimate
	stmt := `SELECT estimate_id, customer_id, created_by, status, created_at,
       	kitchen_length_inch, kitchen_width_inch, kitchen_height_inch,
		door_width_inch, door_height_inch, street, city, state, zip 
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
		var statusInt int
		err = rows.Scan(&estimate.EstimateID, &estimate.CustomerID, &estimate.CreatedBy, &statusInt, &estimate.CreatedAt, &estimate.KitchenLengthInch, &estimate.KitchenWidthInch, &estimate.KitchenHeightInch, &estimate.DoorWidthInch, &estimate.DoorHeightInch, &estimate.Street, &estimate.City, &estimate.State, &estimate.Zip)

		estimate.Status = EstimateStatus(statusInt)

		if err != nil {
			return nil, err
		}

		estimates = append(estimates, estimate)
	}

	return estimates, nil
}

// Update modifies all fields of an existing Estimate record.
// The Estimate struct must contain a valid EstimateID. Returns ErrNoRecord if the record does not exist.
func (m *EstimateModel) Update(e *Estimate) error {
	stmt := `UPDATE estimates
 	SET customer_id=$2, status=$3,
    kitchen_length_inch=$4, kitchen_width_inch=$5, kitchen_height_inch=$6,
    door_width_inch=$7, door_height_inch=$8,
    street=$9, city=$10, state=$11, zip=$12
	WHERE estimate_id=$1`

	result, err := m.DB.Exec(stmt,
		e.EstimateID, e.CustomerID, e.Status,
		e.KitchenLengthInch, e.KitchenWidthInch, e.KitchenHeightInch,
		e.DoorWidthInch, e.DoorHeightInch, e.Street, e.City, e.State, e.Zip,
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

// Delete removes an Estimate from the database using its ID.
// Returns ErrNoRecord if the specified record does not exist.
func (m *EstimateModel) Delete(id int) error {
	stmt := `DELETE FROM estimates WHERE estimate_id=$1`
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

// CalculateEstimateTotals computes the subtotal, labor, sales tax, and total for a given set of EstimateProducts.
// Labor cost logic is based on product categories, and Michigan’s 6% sales tax is applied.
// Returns an EstimateTotals struct with all calculated fields.
func (m *EstimateModel) CalculateEstimateTotals(estimateProducts []EstimateProduct) EstimateTotals {

	var totals EstimateTotals
	for i := 0; i < len(estimateProducts); i++ {
		totals.Subtotal += estimateProducts[i].Product.UnitPrice * estimateProducts[i].EstimateItem.Quantity

		switch estimateProducts[i].Product.Category {
		case "Appliances":
			totals.LaborTotal += 10000
		case "Cabinetry":
			totals.LaborTotal += 2500 * estimateProducts[i].EstimateItem.Quantity
		case "Countertops":
			totals.LaborTotal += 3000 * estimateProducts[i].EstimateItem.Quantity
		case "Flooring":
			totals.LaborTotal += 500 * estimateProducts[i].EstimateItem.Quantity
		case "Sinks & Faucets":
			totals.LaborTotal += 7500 * estimateProducts[i].EstimateItem.Quantity
		}
	}

	totals.LaborTotal += 30000
	totals.SalesTax = totals.Subtotal / 6
	totals.EstimateTotal = totals.Subtotal + totals.SalesTax + totals.LaborTotal

	return totals
}
