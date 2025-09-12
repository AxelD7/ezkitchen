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
    has_island BOOLEAN,
    street VARCHAR(255),
    city VARCHAR(50),
    state VARCHAR(60),
    zip VARCHAR(10)
);
