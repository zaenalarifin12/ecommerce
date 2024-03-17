-- +migrate Up
CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    uuid            VARCHAR(255) UNIQUE NOT NULL,
    full_name       VARCHAR(255)        NOT NULL,
    phone           VARCHAR(10)         NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    username        VARCHAR(50) UNIQUE  NOT NULL,
    password        VARCHAR(100)        NOT NULL,
    email_verify_at TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP
);
