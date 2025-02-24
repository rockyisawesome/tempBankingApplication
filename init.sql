-- -- Create the 'accounts' database
-- CREATE DATABASE accounts;

-- -- Create a user 'postgres' with a password (optional, see note below)
-- CREATE ROLE postgres WITH LOGIN PASSWORD 'abcd';

-- -- Grant privileges to the user on the 'accounts' database
-- GRANT ALL PRIVILEGES ON DATABASE accounts TO postgres;

-- -- Connect to the 'accounts' database
-- \connect accounts

-- Enable pgcrypto extension in the 'accounts' database
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create the 'usersschema' schema
CREATE SCHEMA IF NOT EXISTS usersschema;

-- Create the 'accounts' table in the 'usersschema' schema
CREATE TABLE IF NOT EXISTS usersschema.accounts (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    balance double precision DEFAULT 0.0,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT accounts_pkey PRIMARY KEY (id),
    CONSTRAINT accounts_email_key UNIQUE (email),
    CONSTRAINT accounts_username_key UNIQUE (username)
);

-- Optional: Grant privileges on the schema and table to the user
GRANT USAGE ON SCHEMA usersschema TO postgres;
GRANT ALL PRIVILEGES ON usersschema.accounts TO postgres;