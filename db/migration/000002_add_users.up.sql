create table users (
    username varchar(50) primary key,
    hashed_password varchar not null,
    fullname varchar(50) not null,
    email varchar unique not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    deleted_at timestamptz
);

alter table "accounts" add constraint accounts_owner_fk
    foreign key ("owner") references users(username) on delete cascade;