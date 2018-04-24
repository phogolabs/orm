-- Auto-generated at Tue Apr 24 16:53:27 BST 2018

-- name: select-all-users
SELECT * FROM users

-- name: select-user
SELECT * FROM users
WHERE id = ?

-- name: insert-user
INSERT INTO users (id, first_name, last_name)
VALUES (?, ?, ?)

-- name: update-user
UPDATE users
SET first_name = ?, last_name = ?
WHERE id = ?

-- name: delete-user
DELETE FROM users
WHERE id = ?