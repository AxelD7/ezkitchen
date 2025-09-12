CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(320) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    phone VARCHAR(15),
    created_at TIMESTAMP
);

INSERT INTO users (name, email, password_hash, phone, role)
VALUES (
    'Admin User',
    'admin@example.com',
    'replace_with_hash',  
    '555-123-4567',
    'admin'
);