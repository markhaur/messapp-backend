-- name: GetUserByEmployeeID :one
SELECT * FROM users
WHERE employee_id = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO users(
    name, password, designation, employee_id
) VALUES (
    ?, ?, ?, ?
);

-- name: DeleteUser :exec
DELETE FROM users
where id = ?;

-- name: UpdateUser :exec
UPDATE users SET password = ?
WHERE id = ?;