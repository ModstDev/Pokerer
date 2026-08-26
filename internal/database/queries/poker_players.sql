-- name: AddTablePlayer :exec
INSERT INTO table_players (
    id,
    table_id,
    user_id,
    seat_number,
    chips
)
VALUES (?, ?, ?, ?, ?);

-- name: RemoveTablePlayer :exec
DELETE FROM table_players
WHERE table_id = ?
  AND user_id = ?;

-- name: GetTablePlayer :one
SELECT
    id,
    table_id,
    user_id,
    seat_number,
    chips,
    joined_at
FROM table_players
WHERE table_id = ?
  AND user_id = ?
LIMIT 1;

-- name: ListTablePlayers :many
SELECT
    id,
    table_id,
    user_id,
    seat_number,
    chips,
    joined_at
FROM table_players
WHERE table_id = ?
ORDER BY seat_number;