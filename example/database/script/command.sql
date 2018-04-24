-- name: select-users
SELECT * FROM users

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
