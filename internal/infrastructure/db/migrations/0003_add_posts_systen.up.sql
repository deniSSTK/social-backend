CREATE TABLE IF NOT EXISTS posts (
                                     id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                                     description VARCHAR(1000),
                                     author_id UUID REFERENCES users(id) NOT NULL,
                                     created_at timestamptz NOT NULL DEFAULT NOW(),
                                     likes_count INT NOT NULL DEFAULT 0,
                                     comments_count INT NOT NULL DEFAULT 0,
                                     close_friends BOOLEAN,
                                     pinned BOOLEAN
);

CREATE TABLE IF NOT EXISTS images (
                                      id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                                      url TEXT NOT NULL UNIQUE,
                                      position INT,
                                      delete_url TEXT NOT NULL UNIQUE,
                                      post_id UUID REFERENCES posts(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS post_comments (
                                             id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                                             post_id UUID REFERENCES posts(id) NOT NULL,
                                             author_id UUID REFERENCES users(id) NOT NULL,
                                             created_at timestamptz NOT NULL DEFAULT NOW(),
                                             text VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS post_likes (
                                          post_id UUID REFERENCES posts(id) NOT NULL,
                                          author_id UUID REFERENCES users(id) NOT NULL,
                                          created_at timestamptz NOT NULL DEFAULT NOW(),
                                          PRIMARY KEY (post_id, author_id)
);

CREATE TABLE IF NOT EXISTS hashtags (
                                        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                                        name VARCHAR(100) UNIQUE
);

CREATE TABLE IF NOT EXISTS post_hashtags (
                                             post_id UUID REFERENCES posts(id) NOT NULL,
                                             hashtag_id UUID REFERENCES hashtags(id) NOT NULL,
                                             position INT,
                                             PRIMARY KEY (post_id, hashtag_id)
);