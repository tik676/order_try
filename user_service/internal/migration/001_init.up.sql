CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    password_Hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tokens(
    id BIGSERIAL PRIMARY KEY,
    access_token TEXT UNIQUE NOT NULL,
    refresh_token TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);