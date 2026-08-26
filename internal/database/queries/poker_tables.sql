-- name: CreatePokerTable :exec
INSERT INTO poker_tables (
    id,
    name,
    small_blind,
    big_blind,
    min_buy_in,
    max_buy_in,
    max_players,
    status
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetPokerTableByID :one
SELECT
    id,
    name,
    small_blind,
    big_blind,
    min_buy_in,
    max_buy_in,
    max_players,
    status,
    created_at,
    updated_at
FROM poker_tables
WHERE id = ?
LIMIT 1
FOR UPDATE;

-- name: ListPokerTables :many
SELECT
    pt.id,
    pt.name,
    pt.small_blind,
    pt.big_blind,
    pt.min_buy_in,
    pt.max_buy_in,
    pt.max_players,
    pt.status,
    pt.created_at,
    pt.updated_at,
    COUNT(tp.id) AS player_count
FROM poker_tables pt
LEFT JOIN table_players tp ON tp.table_id = pt.id
GROUP BY
    pt.id,
    pt.name,
    pt.small_blind,
    pt.big_blind,
    pt.min_buy_in,
    pt.max_buy_in,
    pt.max_players,
    pt.status,
    pt.created_at,
    pt.updated_at
ORDER BY pt.created_at DESC;