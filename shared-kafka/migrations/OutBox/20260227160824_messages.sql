-- +goose Up
CREATE TABLE messages (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                          message TEXT NOT NULL,
                          processed BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE messages;