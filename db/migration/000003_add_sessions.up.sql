create table sessions (
    id uuid primary key,
    username varchar(64),
    refresh_token varchar not null,
    user_agent varchar not null,
    client_ip varchar not null,
    is_blocked bool not null default false,
    created_at timestamptz not null default now(),
    deleted_at timestamptz default '0001-01-01 00:00:00Z',
    expires_at timestamptz not null ,
    foreign key (username) references users(username)
);