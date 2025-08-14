-- Initial database schema
-- ---- +migrate Up
-- This migration creates the initial tables for the application

-- Example table (you can customize this according to your needs)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ---- +migrate Down
-- This migration removes the initial tables
DROP TABLE IF EXISTS users;
