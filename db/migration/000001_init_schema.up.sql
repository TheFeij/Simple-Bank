create table accounts(
    id bigserial primary key,
    balance bigint default 0 not null,
    owner varchar(64) not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz default '0001-01-01 00:00:00Z'
);

create table entries(
    id bigserial primary key,
    account_id bigint references accounts(id) on delete cascade not null,
    amount int not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz default '0001-01-01 00:00:00Z'
);

create table transfers(
    id bigserial primary key,
    from_account_id bigint references accounts(id) on delete cascade not null ,
    to_account_id bigint references accounts(id) on delete cascade not null,
    incoming_entry_id bigint references entries(id) on delete cascade not null,
    outgoing_entry_id bigint references entries(id) on delete cascade not null,
    amount int not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    deleted_at timestamptz default '0001-01-01 00:00:00Z'
);

SET timezone to 'UTC'