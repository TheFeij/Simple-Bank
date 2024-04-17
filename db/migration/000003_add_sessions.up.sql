create table sessions (
    id uuid primary key,
    username varchar(64),
    refresh_token varchar not null,
    user_agent varchar not null,
    client_ip varchar unique not null,
    is_blocked bool not null default false,
    created_at timestamptz not null default now(),
    expires_at timestamptz not null ,
    foreign key (username) references users(username)
);