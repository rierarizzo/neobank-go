CREATE SCHEMA IF NOT EXISTS account AUTHORIZATION CURRENT_USER;
SET search_path TO account;

CREATE TABLE customers (
    id              UUID PRIMARY KEY,
    identity_number VARCHAR(100) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL,
    phone_number    VARCHAR(100) NOT NULL,
    date_of_birth   DATE,
    nationality     VARCHAR(50),
    address_line1   VARCHAR(255),
    address_line2   VARCHAR(255),
    city            VARCHAR(100),
    state           VARCHAR(100),
    postal_code     VARCHAR(20),
    country         VARCHAR(100),
    created_at      timestamptz  NOT NULL DEFAULT NOW(),
    updated_at      timestamptz  NOT NULL DEFAULT NOW());

CREATE INDEX idx_customers_identity_number ON customers (identity_number);

CREATE TABLE accounts (
    id                UUID PRIMARY KEY,
    customer_id       UUID        NOT NULL,
    ledger_account_id BIGSERIAL   NOT NULL,
    account_type      VARCHAR(50) NOT NULL DEFAULT 'checking',
    status            VARCHAR(15) NOT NULL DEFAULT 'active',
    opened_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at         TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_accounts_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON UPDATE CASCADE ON DELETE RESTRICT)