-- name: CreateWallet :exec
INSERT INTO wallets (
    id,
    user_id
)
VALUES (?, ?);

-- name: GetWalletByUserID :one
SELECT
    id,
    user_id,
    balance,
    created_at,
    updated_at
FROM wallets
WHERE user_id = ?
LIMIT 1;

-- name: CreateWalletTransaction :exec
INSERT INTO wallet_transactions (
    id,
    wallet_id,
    type,
    amount,
    balance_after
)
VALUES (?, ?, ?, ?, ?);

-- name: GetWalletTransactionsByUserID :many
SELECT
    wt.id,
    wt.wallet_id,
    wt.type,
    wt.amount,
    wt.balance_after,
    wt.created_at
FROM wallet_transactions wt
JOIN wallets w ON w.id = wt.wallet_id
WHERE w.user_id = ?
ORDER BY wt.created_at DESC;