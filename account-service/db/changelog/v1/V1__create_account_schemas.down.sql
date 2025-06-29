SET search_path TO account;

DROP INDEX IF EXISTS idx_customers_identity_number;

DROP TABLE IF EXISTS account.accounts;
DROP TABLE IF EXISTS account.customers;

DROP SCHEMA IF EXISTS account;
