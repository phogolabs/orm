-- Auto-generated at Sat, 16 Mar 2019 21:43:37 CET
-- Please do not change the name attributes

-- name: up
CREATE TABLE IF NOT EXISTS migrations (
 id          VARCHAR(15) NOT NULL PRIMARY KEY,
 description TEXT        NOT NULL,
 created_at  TIMESTAMP   NOT NULL
);

-- name: down
DROP TABLE IF EXISTS migrations;
