package integration_test

import (
	"context"
	"database/sql"
	"ezkitchen/internal/models"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testDB            *sql.DB
	userModel         *models.UserModel
	estimateModel     *models.EstimateModel
	productModel      *models.ProductModel
	estimateItemModel *models.EstimateItemModel
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine", postgres.WithDatabase("ezkitchen_test"), postgres.WithUsername("testuser"), postgres.WithPassword("testpass"), postgres.BasicWaitStrategies())
	if err != nil {
		log.Fatalf("Failed to start postgres container %v", err)
	}

	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Printf("Could not terminate container: %v", err)

		}

	}()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err := createSchema(db); err != nil {
		log.Fatalf("Schema init failed: %v", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}

	testDB = db
	userModel = &models.UserModel{DB: db}
	estimateModel = &models.EstimateModel{DB: db}
	productModel = &models.ProductModel{DB: db}
	estimateItemModel = &models.EstimateItemModel{DB: db}

	code := m.Run()

	_ = db.Close()
	os.Exit(code)

}

func createSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(320) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    phone VARCHAR(15),
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    category VARCHAR(50),
    subcategory VARCHAR(50),
    color VARCHAR(20),
    unit_price INT NOT NULL,
    length REAL,
    width REAL,
    height REAL,
    created_by INT REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS estimates (
    estimate_id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES users(user_id),
    created_by INT NOT NULL REFERENCES users(user_id),
    status INT CHECK (status >= 1 AND status <= 6),
    created_at TIMESTAMP,
    kitchen_length_inch DOUBLE PRECISION,
    kitchen_width_inch DOUBLE PRECISION,
    kitchen_height_inch DOUBLE PRECISION,
    door_width_inch DOUBLE PRECISION,
    door_height_inch DOUBLE PRECISION,
    flooring_type VARCHAR(255),
    street VARCHAR(255),
    city VARCHAR(50),
    state VARCHAR(60),
    zip VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS estimate_items (
    line_item_id SERIAL PRIMARY KEY,
    estimate_id INT NOT NULL REFERENCES estimates(estimate_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id),
    quantity INT NOT NULL DEFAULT 1
);
`
	_, err := db.Exec(schema)
	return err
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`TRUNCATE estimate_items, estimates, products, users RESTART IDENTITY CASCADE;`)
	if err != nil {
		t.Fatalf("resetDB failed: %v", err)
	}
}
