create table if not exists persons (
    id bigint generated always as identity primary key,
    name text not null,
    age int,
    address text not null,
    work text not null
);