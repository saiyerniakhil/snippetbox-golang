-- name: GetSnippetById :one
SELECT id, title, content, created, expires
FROM snippets
WHERE expires > UTC_TIMESTAMP()
  and id = ?;

-- name: GetLatestSnippets :many
SELECT id, title, content, created, expires
from snippets
WHERE expires > UTC_TIMESTAMP()
order by id DESC
LIMIT ?;

-- name: AddSnippet :execresult
INSERT INTO snippets (title, content, expires, created)
VALUES (?, ?, ?, UTC_TIMESTAMP());

-- name: AddUser :execresult
INSERT INTO users (name, email, hashed_password, created)
VALUES (?, ?, ?, UTC_TIMESTAMP());

-- name: AuthenticateUser :one
SELECT id, hashed_password
FROM users
WHERE email = ?;

-- name: UserExists :one
SELECT EXISTS(SELECT true FROM users WHERE id = ?);