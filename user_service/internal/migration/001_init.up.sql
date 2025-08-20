CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    password_Hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tokens(
    id BIGSERIAL PRIMARY KEY,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    user_id BIGINT REFERENCES users(id)
);