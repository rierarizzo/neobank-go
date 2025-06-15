set search_path to ledger;

drop index if exists idx_entries_account_id;
drop index if exists idx_entries_transfer_id;

drop table if exists ledger.entries;
drop table if exists ledger.ledger_accounts;
drop table if exists ledger.transfers;

drop
    type if exists account_status;
drop
    domain if exists currency;

-- (Si creaste el schema y quieres limpiarlo)
drop schema if exists ledger;
