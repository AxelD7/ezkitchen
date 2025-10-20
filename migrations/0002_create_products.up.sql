CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(255),
    category VARCHAR(50),
    subcategory VARCHAR(50),
    color        VARCHAR(20),
    unit_price INT,
    created_by INT REFERENCES users(user_id)
);
