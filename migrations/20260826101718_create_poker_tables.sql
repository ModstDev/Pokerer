-- +goose Up

CREATE TABLE poker_tables (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    small_blind BIGINT NOT NULL,
    big_blind BIGINT NOT NULL,
    min_buy_in BIGINT NOT NULL,
    max_buy_in BIGINT NOT NULL,
    max_players INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE table_players (
    id CHAR(36) NOT NULL PRIMARY KEY,
    table_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    seat_number INT NOT NULL,
    chips BIGINT NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_table_players_table
        FOREIGN KEY (table_id) REFERENCES poker_tables(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_table_players_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_table_player
        UNIQUE (table_id, user_id),

    CONSTRAINT uq_table_seat
        UNIQUE (table_id, seat_number),

    INDEX idx_table_player_user (user_id)
);

-- +goose Down

DROP TABLE table_players;
DROP TABLE poker_tables;
