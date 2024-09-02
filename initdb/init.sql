-- Создание базы данных face_scanner
CREATE DATABASE face_scanner;

-- Подключение к базе данных face_scanner
\c face_scanner

CREATE TABLE IF NOT EXISTS image_tasks
(
    task_id      TEXT DEFAULT ''::text NOT NULL,
    image_data   BYTEA,
    api_response TEXT,
    image_id     TEXT DEFAULT ''::text NOT NULL,
    file_name    TEXT DEFAULT ''::text NOT NULL
);

ALTER TABLE image_tasks
    OWNER TO root;

CREATE TABLE IF NOT EXISTS tasks
(
    id     TEXT    DEFAULT ''::text NOT NULL UNIQUE,
    status INTEGER DEFAULT 1        NOT NULL
);

ALTER TABLE tasks
    OWNER TO root;

CREATE TABLE IF NOT EXISTS user_data
(
    email     TEXT    DEFAULT ''::text NOT NULL UNIQUE,
    hash_password TEXT DEFAULT ''::text       NOT NULL
);

ALTER TABLE user_data
    OWNER TO root;