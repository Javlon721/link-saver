DROP table if EXISTS users CASCADE;

create table
  if not exists users (
    id serial primary key,
    telegram_id bigint not null unique
  );