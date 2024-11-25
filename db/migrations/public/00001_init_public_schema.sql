-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA pg_catalog;

-- Create files table
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP DEFAULT NULL
);


-- Create the tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo UUID,
    workspace VARCHAR(255) NOT NULL UNIQUE DEFAULT md5(random()::text),  -- default random string
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (logo) REFERENCES files(id)
);

-- Create user_lookup table
CREATE TABLE IF NOT EXISTS user_lookup (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Add check constraint for email format, only when email is not null
ALTER TABLE user_lookup ADD CONSTRAINT check_email_format
  CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables with cascading effect
DROP TABLE IF EXISTS user_lookup CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;
DROP TABLE IF EXISTS files CASCADE;

-- +goose StatementEnd
