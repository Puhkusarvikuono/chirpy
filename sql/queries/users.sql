-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: Reset :exec
DELETE FROM users;

<<<<<<< HEAD

=======
-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 limit 1;
>>>>>>> eae180b (feat(auth): implement HS256 JWT signing and refresh token persistence in Go.)
