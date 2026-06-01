-- name: GetUserByName :one
select * from users
where name = $1;

-- name: GetUserByID :one
select * from users
where id = $1;

-- name: GetUserByClerkUserID :one
select * from users
where clerk_user_id = $1;

-- name: CreateUser :one
insert into users (
  clerk_user_id,
  name,
  display_name,
  avatar_url,
  email
) values (
  $1,
  $2,
  $3,
  $4,
  $5
)
returning *;

-- name: UpdateUserProfile :one
update users
set
  name = $2,
  display_name = $3,
  avatar_url = $4,
  description = $5
where id = $1
returning *;

-- name: DeleteUser :exec
delete from users
where id = $1;
