alter table repositories
    add column item_prefix      text   not null default 'I',
    add column next_item_number bigint not null default 1;

create type repository_item_kind as enum ('issue', 'pull_request');
create type issue_status as enum ('open', 'in_progress', 'completed', 'cancelled');

create table repository_items
(
    id               uuid primary key                  default uuidv7(),
    repository_id    uuid                     not null references repositories (id) on delete cascade,
    number           bigint                   not null,
    kind             repository_item_kind     not null,
    title            text                     not null,
    body             text,
    author_id        uuid                     references users (id) on delete set null,
    created_at       timestamp with time zone not null default now(),
    updated_at       timestamp with time zone not null default now(),
    last_activity_at timestamp with time zone not null default now(),

    constraint repository_items_repository_number_unique unique (repository_id, number),
    constraint repository_items_id_kind_unique unique (id, kind)
);

create trigger repository_items_set_updated_at
    before update
    on repository_items
    for each row
execute function set_updated_at();

create table issues
(
    repository_item_id uuid primary key,
    kind               repository_item_kind not null default 'issue', -- for item kind constraint, should never be touched by application code
    status             issue_status         not null default 'open',

    constraint issues_kind_issue check (kind = 'issue'),
    constraint issues_repository_item_fk
        foreign key (repository_item_id, kind)
            references repository_items (id, kind)
            on delete cascade
);

create table labels
(
    id            uuid primary key                  default uuidv7(),
    repository_id uuid                     not null references repositories (id) on delete cascade,
    name          text                     not null,
    color         text                     not null,
    description   text,
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now(),

    constraint labels_name_lowercase_check check (name = lower(name)),
    constraint labels_repository_name_unique unique (repository_id, name)
);

create trigger labels_set_updated_at
    before update
    on labels
    for each row
execute function set_updated_at();

create table repository_item_labels
(
    repository_item_id uuid                     not null references repository_items (id) on delete cascade,
    label_id           uuid                     not null references labels (id) on delete cascade,
    created_at         timestamp with time zone not null default now(),

    primary key (repository_item_id, label_id)
);

create table repository_item_timeline
(
    id                 uuid primary key                  default uuidv7(),
    repository_item_id uuid                     not null references repository_items (id) on delete cascade,
    actor_id           uuid                     references users (id) on delete set null,
    payload            jsonb                    not null,
    created_at         timestamp with time zone not null default now(),
    updated_at         timestamp with time zone not null default now()
);

create trigger repository_item_timeline_set_updated_at
    before update
    on repository_item_timeline
    for each row
execute function set_updated_at();

create index repository_items_repository_last_activity_at_idx on repository_items (repository_id, last_activity_at desc);
create index repository_items_author_id_idx on repository_items (author_id);
create index labels_repository_id_idx on labels (repository_id);
create index repository_item_labels_label_id_idx on repository_item_labels (label_id);
create index repository_item_timeline_item_created_at_idx on repository_item_timeline (repository_item_id, created_at);

---- create above / drop below ----

drop table repository_item_timeline;
drop table repository_item_labels;
drop table labels;
drop table issues;
drop table repository_items;
drop type issue_status;
drop type repository_item_kind;
alter table repositories
    drop column next_item_number,
    drop column item_prefix;
