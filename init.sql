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

DROP TABLE IF EXISTS usersschema.accounts;
DROP TABLE IF EXISTS usersschema.transactions;

-- Create the 'accounts' table in the 'usersschema' schema
CREATE TABLE usersschema.accounts (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    account_number character varying(255) NOT NULL, -- Unique account number
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default" NOT NULL,
    balance double precision DEFAULT 0.0,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT accounts_pkey PRIMARY KEY (id),
    CONSTRAINT accounts_accountnumber_key UNIQUE (account_number)
);


-- Create the transactions table
CREATE TABLE usersschema.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- UUID primary key with default generation
    from_account_id character varying(255) NOT NULL,
    to_account_id character varying(255), -- Nullable for deposits/withdrawals involving external systems
    amount double precision NOT NULL CHECK (amount > 0), -- Precision for currency (e.g., 15 digits, 2 after decimal)
    transaction_type VARCHAR(50) NOT NULL CHECK (transaction_type IN ('transfer', 'deposit', 'withdrawal')),
    description TEXT, -- Optional field, can store longer text
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')),
    -- Foreign key constraints linking to accounts table
    CONSTRAINT fk_from_account FOREIGN KEY (from_account_id) REFERENCES usersschema.accounts(account_number) ON DELETE RESTRICT,
    CONSTRAINT fk_to_account FOREIGN KEY (to_account_id) REFERENCES usersschema.accounts(account_number) ON DELETE RESTRICT
);

-- Optional: Grant privileges on the schema and table to the user
GRANT USAGE ON SCHEMA usersschema TO postgres;
GRANT ALL PRIVILEGES ON usersschema.accounts TO postgres;
GRANT ALL PRIVILEGES ON usersschema.transactions TO postgres;