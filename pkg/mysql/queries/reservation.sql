-- name: GetReservationsByEmployeeID :many
SELECT reservations.*, users.name FROM reservations
INNER JOIN users ON reservations.user_id = users.id
WHERE user_id = ?;

-- name: GetReservationByID :one
SELECT reservations.*, users.name FROM reservations
INNER JOIN users ON reservations.user_id = users.id
WHERE reservations.id = ? LIMIT 1;

-- name: GetReservationsByDate :many
SELECT reservations.*, users.name FROM reservations
INNER JOIN users ON reservations.user_id = users.id
WHERE DATE(reservation_time) = DATE(?);

-- name: ListReservations :many
SELECT reservations.*, users.name FROM reservations
INNER JOIN users ON reservations.user_id = users.id
ORDER BY reservation_time desc;

-- name: CreateReservation :execresult
INSERT INTO reservations (
    user_id, reservation_time, type, no_of_guests, created_at
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: DeleteReservation :exec
DELETE FROM reservations
where id = ?;

-- name: UpdateReservation :exec
UPDATE reservations SET no_of_guests = ?
WHERE id = ?;