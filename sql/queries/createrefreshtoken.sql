-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5);


-- name: GetUserFromRefreshToken :one
SELECT users.*
FROM users
JOIN refresh_tokens ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
	AND refresh_tokens.expires_at > now()
	AND refresh_tokens.revoked_at is NULL;


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = now(),
updated_at = now()
WHERE token = $1;
