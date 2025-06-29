SET search_path TO ledger;

DROP INDEX IF EXISTS idx_entries_account_id;
DROP INDEX IF EXISTS idx_entries_transfer_id;

DROP TABLE IF EXISTS ledger.entries;
DROP TABLE IF EXISTS ledger.ledger_accounts;
DROP TABLE IF EXISTS ledger.transfers;

-- (Si creaste el schema y quieres limpiarlo)
DROP SCHEMA IF EXISTS ledger;
