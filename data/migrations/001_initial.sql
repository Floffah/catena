create table users (
  id serial primary key,
  clerk_user_id text unique not null,
  name text not null unique,
  created_at timestamp with time zone default now(),
  updated_at timestamp with time zone default now()
);
create index on users (name);

---- create above / drop below ----

drop table users;