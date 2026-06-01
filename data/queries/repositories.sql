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

-- name: ListRepositoriesByOwnerUpdated :many
select *
from repositories
where owner_id = sqlc.arg(owner_id)
  and (
    (
      not sqlc.arg(filter_visibility)::boolean
      and (
        visibility = 'public'
        or sqlc.arg(include_private)::boolean
      )
    )
    or (
      sqlc.arg(filter_visibility)::boolean
      and visibility = sqlc.arg(visibility)::repository_visibility
      and (
        visibility = 'public'
        or sqlc.arg(include_private)::boolean
      )
    )
  )
order by updated_at desc, name asc
limit sqlc.arg(result_limit);

-- name: ListRepositoriesByOwnerFeatured :many
select *
from repositories
where owner_id = sqlc.arg(owner_id)
  and (
    (
      not sqlc.arg(filter_visibility)::boolean
      and (
        visibility = 'public'
        or sqlc.arg(include_private)::boolean
      )
    )
    or (
      sqlc.arg(filter_visibility)::boolean
      and visibility = sqlc.arg(visibility)::repository_visibility
      and (
        visibility = 'public'
        or sqlc.arg(include_private)::boolean
      )
    )
  )
order by updated_at desc, name asc
limit sqlc.arg(result_limit);

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
