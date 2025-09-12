CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(255),
    unit_price NUMERIC(12,2),
    created_by INT REFERENCES users(user_id)
);
