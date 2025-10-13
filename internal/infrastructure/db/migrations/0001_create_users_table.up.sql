DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
            CREATE TYPE user_status AS ENUM ('ACTIVE', 'BLOCK');
        END IF;
    END$$;

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     username VARCHAR(50) NOT NULL,
                                     email VARCHAR(100) UNIQUE,
                                     password_hash TEXT,
                                     status user_status DEFAULT 'ACTIVE',
                                     icon_url TEXT,
                                     description VARCHAR(100),
                                     created_at timestamptz NOT NULL DEFAULT NOW(),
                                     updated_at timestamptz NOT NULL DEFAULT NOW()
)