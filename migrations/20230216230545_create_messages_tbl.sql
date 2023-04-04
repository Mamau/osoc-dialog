-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    text TEXT,
    user_id  INT NOT NULL,
    author_id  INT NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
