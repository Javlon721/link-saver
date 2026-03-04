DROP table if EXISTS links CASCADE;

create table
  if not exists links (
    id serial primary key,
    user_id bigint references users (id) not null,
    link varchar not null,
    describtion varchar
  );