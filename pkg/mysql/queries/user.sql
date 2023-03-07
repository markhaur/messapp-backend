--name: GetUser :one
SELECT * FROM users
WHERE employee_id = ? LIMIT 1;

--name: ListUsers :many
SELECT * FROM users
ORDER BY name;

--name: CreateUser :one
INSERT INTO users(
    name, password, designation, employee_id
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

--name: DeleteUser :exec
DELETE FROM users
where id = ?;

--name: UpdateUser :exec
UPDATE users SET password = ?
WHERE id = ?;