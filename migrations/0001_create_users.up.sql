CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(60) NULL,
    role VARCHAR(20) NOT NULL,
    phone VARCHAR(15),
    created_at TIMESTAMP
);

INSERT INTO users (name, email, hashed_password, phone, role)
VALUES (
    'Admin User',
    'admin@example.com',
    'replace_with_hash',  
    '555-123-4567',
    'admin'
);