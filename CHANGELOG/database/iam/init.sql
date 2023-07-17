CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    domain VARCHAR(255),
    name VARCHAR(255),
    created_at BIGINT,
    updated_at BIGINT,
    hashed_password VARCHAR(255),
    display_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255),
    gender VARCHAR(255),
    is_admin BOOLEAN,
    is_forbidden BOOLEAN,
    CONSTRAINT domain_name UNIQUE (domain, name)
);
