-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                     telegram_id BIGINT UNIQUE,
                                     email VARCHAR(255),
    username VARCHAR(64),
    credits INT DEFAULT 0,
    is_blocked BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
                                            id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                            message_id VARCHAR(64),
    user_id BIGINT,
    email VARCHAR(255),
    type VARCHAR(32),
    tier_name VARCHAR(64),
    amount DECIMAL(10,2),
    currency VARCHAR(8),
    data TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS grpc_queue (
                                          id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                          user_id BIGINT,
                                          order_id VARCHAR(255),
    amount INT,
    credits INT,
    email VARCHAR(255),
    username VARCHAR(64),
    provider VARCHAR(64),
    attempts INT DEFAULT 0,
    last_error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS grpc_queue;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
