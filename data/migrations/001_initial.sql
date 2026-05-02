create function set_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

create type repository_visibility as enum ('private', 'public');

create table users (
  id uuid primary key default uuidv7(),
  clerk_user_id text not null unique,
  name text not null unique,
  display_name text,
  avatar_url text,
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  constraint users_name_not_empty check (length(trim(name)) > 0)
);

create trigger users_set_updated_at
before update on users
for each row
execute function set_updated_at();

create table repositories (
  id uuid primary key default uuidv7(),
  owner_id uuid not null references users (id) on delete cascade,
  name text not null,
  description text,
  visibility repository_visibility not null default 'private',
  default_branch text not null default 'main',
  created_at timestamp with time zone not null default now(),
  updated_at timestamp with time zone not null default now(),

  constraint repositories_owner_name_unique unique (owner_id, name)
);

create trigger repositories_set_updated_at
before update on repositories
for each row
execute function set_updated_at();

create index users_clerk_user_id_idx on users (clerk_user_id);
create index repositories_owner_id_idx on repositories (owner_id);
create index repositories_visibility_idx on repositories (visibility);
create index repositories_owner_updated_at_idx on repositories (owner_id, updated_at desc);

---- create above / drop below ----

drop table repositories;
drop table users;
drop type repository_visibility;
drop function set_updated_at();
