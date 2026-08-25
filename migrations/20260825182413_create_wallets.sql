-- +goose Up
CREATE TABLE wallets (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36) NOT NULL UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_wallets_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE wallet_transactions (
    id CHAR(36) NOT NULL PRIMARY KEY,
    wallet_id CHAR(36) NOT NULL,
    type VARCHAR(30) NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_wallet_transactions_wallet
        FOREIGN KEY (wallet_id) REFERENCES wallets(id),

    INDEX idx_wallet_transactions_wallet_created (
        wallet_id,
        created_at
    )
);

-- +goose Down
DROP TABLE wallet_transactions;
DROP TABLE wallets;