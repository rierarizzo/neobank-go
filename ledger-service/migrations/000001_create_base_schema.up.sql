create schema if not exists ledger authorization current_user;
set search_path to ledger;

create domain currency as char(3)
    check (value in ('USD', 'EUR'));

create type account_status as enum ('ACTIVE','FROZEN','CLOSED');

-- Tabla de cuentas, con UNIQUE(id, currency) para la FK compuesta
create table ledger_accounts (
    id       bigserial primary key,
    currency currency       not null,
    status   account_status not null default 'ACTIVE',
    constraint uq_accounts_id_currency unique (id, currency));

-- Tabla de transferencias, referenciando (id, currency) de ledger_accounts
create table transfers (
    id              bigserial primary key,
    external_id     text        not null unique,
    from_account_id bigint      not null,
    to_account_id   bigint      not null,
    currency        currency    not null,
    created_at      timestamptz not null default now(),
    constraint fk_transfers_from
        foreign key (from_account_id, currency)
            references ledger_accounts (id, currency)
            on delete restrict,
    constraint fk_transfers_to
        foreign key (to_account_id, currency)
            references ledger_accounts (id, currency)
            on delete restrict);

create table entries (
    id                bigserial primary key,
    transfer_id       bigint not null,
    ledger_account_id bigint not null,
    amount            int    not null,
    foreign key (transfer_id) references transfers (id) on delete restrict,
    foreign key (ledger_account_id) references ledger_accounts (id) on delete restrict,
    check ( amount <> 0 ));

create index idx_entries_transfer_id on entries (transfer_id);
create index idx_entries_account_id on entries (ledger_account_id);
