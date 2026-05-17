-- name: CreateGitAccessToken :one
insert into git_access_tokens (
  user_id,
  name,
  token_hash,
  token_prefix,
  scopes,
  expires_at
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
returning *;

-- name: ListGitAccessTokensByUserID :many
select * from git_access_tokens
where user_id = $1
  and revoked_at is null
  and (expires_at is null or expires_at > now())
order by created_at desc;

-- name: GetGitAccessTokenByHash :one
select * from git_access_tokens
where token_hash = $1;

-- name: RevokeGitAccessToken :exec
update git_access_tokens
set revoked_at = now()
where id = $1
  and user_id = $2
  and revoked_at is null;

-- name: TouchGitAccessTokenLastUsed :exec
update git_access_tokens
set last_used_at = now()
where id = $1;
