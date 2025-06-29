CREATE SCHEMA IF NOT EXISTS ledger AUTHORIZATION CURRENT_USER;
SET search_path TO ledger;

-- Tabla de cuentas
CREATE TABLE ledger_accounts (
    id         BIGSERIAL PRIMARY KEY,
    currency   CHAR(3)     NOT NULL,
    status     VARCHAR(10) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (id, currency));

-- Tabla de transferencias
CREATE TABLE transfers (
    id              BIGSERIAL PRIMARY KEY,
    external_id     TEXT        NOT NULL UNIQUE,
    from_account_id BIGSERIAL   NOT NULL,
    to_account_id   BIGSERIAL   NOT NULL,
    currency        CHAR(3)     NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (from_account_id, currency)
        REFERENCES ledger_accounts (id, currency) ON DELETE RESTRICT,
    FOREIGN KEY (to_account_id, currency)
        REFERENCES ledger_accounts (id, currency) ON DELETE RESTRICT);

-- Tabla de asientos
CREATE TABLE entries (
    id                BIGSERIAL PRIMARY KEY,
    transfer_id       BIGSERIAL NOT NULL REFERENCES transfers (id) ON DELETE RESTRICT,
    ledger_account_id BIGSERIAL NOT NULL REFERENCES ledger_accounts (id) ON DELETE RESTRICT,
    amount            INT       NOT NULL CHECK (amount <> 0));

-- Índices
CREATE INDEX idx_entries_transfer_id ON entries (transfer_id);
CREATE INDEX idx_entries_account_id ON entries (ledger_account_id);
