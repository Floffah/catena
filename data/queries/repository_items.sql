-- name: ReserveRepositoryItemNumber :one
update repositories
set next_item_number = next_item_number + 1
where id = $1
returning (next_item_number - 1)::bigint as number;

-- name: CreateRepositoryItem :one
insert into repository_items (
  repository_id,
  number,
  kind,
  title,
  body,
  author_id
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
returning *;
