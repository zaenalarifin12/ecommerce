-- +migrate Up
CREATE TABLE IF NOT EXISTS carts
(
    id           SERIAL PRIMARY KEY,
    uuid         VARCHAR(255) UNIQUE NOT NULL,
    product_uuid VARCHAR(255)        NOT NULL,
    user_uuid    VARCHAR(255)         NOT NULL,
    quantity     INT       DEFAULT 1,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP,
    deleted_at   TIMESTAMP
);
