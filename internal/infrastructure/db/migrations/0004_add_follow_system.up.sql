ALTER TABLE users
    ADD COLUMN IF NOT EXISTS post_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS followers INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS following INT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS followings (
                                          follower_id UUID NOT NULL REFERENCES users(id),
                                          follow_to_id UUID NOT NULL REFERENCES users(id),
                                          follow_at timestamptz NOT NULL DEFAULT NOW(),
                                          PRIMARY KEY (follower_id, follow_to_id)
);