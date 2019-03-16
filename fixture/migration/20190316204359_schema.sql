-- Auto-generated at Sat, 16 Mar 2019 20:43:59 UTC
-- Please do not change the name attributes

-- name: up
CREATE TABLE users (
  id         TEXT PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name  TEXT NOT NULL
);

-- name: down
DROP TABLE IF EXISTS users;

