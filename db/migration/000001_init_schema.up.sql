create table accounts(
    id bigserial primary key,
    balance bigint default 0 not null,
    owner varchar(50) not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    deleted_at timestamptz
);

create table entries(
    id bigserial primary key,
    account_id bigint references accounts(id) on delete cascade not null,
    amount bigint not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    deleted_at timestamptz
);

create table transfers(
    id bigserial primary key,
    from_account_id bigint references accounts(id) on delete cascade not null ,
    to_account_id bigint references accounts(id) on delete cascade not null,
    incoming_entry_id bigint references entries(id) on delete cascade not null,
    outgoing_entry_id bigint references entries(id) on delete cascade not null,
    amount bigint not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    deleted_at timestamptz
)