-- name: CreateRepository :one
insert into repositories (
  owner_id,
  name,
  description,
  visibility,
  default_branch
) values (
  $1,
  $2,
  $3,
  $4,
  $5
)
returning *;

-- name: GetRepositoryByID :one
select * from repositories
where id = $1;

-- name: GetRepositoryByOwnerAndName :one
select repositories.*
from repositories
join users on users.id = repositories.owner_id
where users.name = sqlc.arg(owner_name)
  and repositories.name = sqlc.arg(repository_name);

-- name: UpdateRepository :one
update repositories
set
  name = $2,
  description = $3,
  visibility = $4,
  default_branch = $5
where id = $1
returning *;

-- name: DeleteRepository :exec
delete from repositories
where id = $1;
