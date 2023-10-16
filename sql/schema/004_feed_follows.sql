-- +goose Up
CREATE TABLE feed_follow (
    ID UUID PRIMARY KEY,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE (feed_id, user_id)
);
-- +goose Down
DROP TABLE feed_follow;
