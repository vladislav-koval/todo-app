CREATE SCHEMA IF NOT EXISTS todoapp;

CREATE TABLE IF NOT EXISTS todoapp.users
(
    id           SERIAL PRIMARY KEY,
    version      INT          NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_length(trim(full_name)) BETWEEN 3 AND 100),
    phone_number VARCHAR(15) CHECK (
        phone_number ~ '^\+[0-9]+$'
            AND
        char_length(phone_number) BETWEEN 10 AND 15
        ),
    created_at   TIMESTAMPTZ           DEFAULT now()
);

CREATE TABLE IF NOT EXISTS todoapp.tasks
(
    id             SERIAL PRIMARY KEY,
    version        INT          NOT NULL DEFAULT 1,
    title          VARCHAR(100) NOT NULL CHECK (char_length(trim(title)) BETWEEN 1 AND 100),
    description    VARCHAR(1000) CHECK (char_length(trim(description)) BETWEEN 1 AND 1000),
    completed      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    completed_at   TIMESTAMPTZ,
    author_user_id INTEGER      NOT NULL REFERENCES todoapp.users (id),

    CHECK (
        (completed = FALSE AND completed_at IS NULL)
            OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
        )
);