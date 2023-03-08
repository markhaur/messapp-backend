-- name: GetReservationsByEmployeeID :many
SELECT * FROM reservations
WHERE user_id = ?;

-- name: GetReservationByID :one
SELECT * FROM reservations
WHERE id = ? LIMIT 1;

-- name: GetReservationsByDate :many
SELECT * FROM reservations
where reservation_time = ?;

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