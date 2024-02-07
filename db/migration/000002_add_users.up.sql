create table users (
    username varchar(64) primary key,
    hashed_password varchar not null,
    fullname varchar(64) not null,
    email varchar unique not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    deleted_at timestamptz default '0001-01-01 00:00:00Z' not null
);

alter table "accounts" add constraint accounts_owner_fk
    foreign key ("owner") references users(username) on delete cascade;