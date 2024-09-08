create table if not exists public.users(
    id bigint generated always as identity primary key,
    name text unique not null
);