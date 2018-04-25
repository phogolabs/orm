-- Auto-generated at Wed Apr 25 11:26:24 BST 2018

-- name: select-all-users
SELECT * FROM users;

-- name: select-user
SELECT * FROM users
WHERE id = ?;

-- name: insert-user
INSERT INTO users (id, first_name, last_name)
VALUES (?, ?, ?);

-- name: update-user
UPDATE users
SET first_name = ?, last_name = ?
WHERE id = ?;

-- name: delete-user
DELETE FROM users
WHERE id = ?;
