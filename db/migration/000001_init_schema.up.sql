create table accounts(
                         id bigserial primary key,
                         balance bigint default 0 not null,
                         owner varchar(50) not null,
                         createdAt timestamptz default now() not null,
                         currency varchar(15) not null
);

create table entries(
                        id bigserial primary key,
                        account_id bigint references accounts(id) not null,
                        amount bigint not null,
                        createdAt timestamptz default now() not null
);

create table transfers(
                          id bigserial primary key,
                          from_account_id bigint references accounts(id) not null,
                          to_account_id bigint references accounts(id) not null,
                          amount bigint not null,
                          createdAt timestamptz default now() not null
)