CREATE TABLE IF NOT EXISTS users
(
    user_id  SERIAL PRIMARY KEY,
    user_name  VARCHAR(255) NOT NULL,
    surname  VARCHAR(255)        NOT NULL,
    email    VARCHAR(255) UNIQUE NOT NULL,
    hashed_phone    VARCHAR(255)        NOT NULL,
    salt            int NOT NULL,
    domain_number  int   NOT NULL CHECK (domain_number >= 1 AND domain_number <= 16),
    created_at     timestamp NOT NULL
)

