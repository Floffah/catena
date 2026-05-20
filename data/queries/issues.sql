-- name: CreateIssue :one
insert into issues (
  repository_item_id
) values (
  $1
)
returning *;

-- name: GetIssueByRepositoryAndNumber :one
select
  repository_items.id,
  repository_items.repository_id,
  repository_items.number,
  repository_items.kind,
  repository_items.title,
  repository_items.body,
  repository_items.author_id,
  repository_items.created_at,
  repository_items.updated_at,
  repository_items.last_activity_at,
  issues.status
from repository_items
join issues on issues.repository_item_id = repository_items.id
where repository_items.repository_id = sqlc.arg(repository_id)
  and repository_items.number = sqlc.arg(number);

-- name: ListIssuesByRepository :many
select
  repository_items.id,
  repository_items.repository_id,
  repository_items.number,
  repository_items.kind,
  repository_items.title,
  repository_items.body,
  repository_items.author_id,
  repository_items.created_at,
  repository_items.updated_at,
  repository_items.last_activity_at,
  issues.status
from repository_items
join issues on issues.repository_item_id = repository_items.id
where repository_items.repository_id = $1
order by repository_items.last_activity_at desc, repository_items.number desc;
