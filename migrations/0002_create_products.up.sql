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
