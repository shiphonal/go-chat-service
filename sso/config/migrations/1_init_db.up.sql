CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    role INTEGER DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS permission
(
    id     INTEGER PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    "get"    BOOLEAN,
    "update"    BOOLEAN
);

INSERT INTO permission (name, "get", "update")
VALUES ('user', true, true),
       ('mod', true, true),
       ('admin', true, true),
       ('ban', false, false)