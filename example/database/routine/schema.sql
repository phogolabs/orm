-- Auto-generated at Wed Apr 25 11:26:24 BST 2018

-- name: select-all-users
SELECT * FROM users;

-- name: select-user
SELECT * FROM users
WHERE id = ?;

-- name: insert-user
INSERT INTO users (id, first_name, last_name)
VALUES (?, ?, ?);

-- name: search-user-by-name
SELECT * FROM users
{{#if name}}
WHERE first_name LIKE '%{{name}}%' OR last_name LIKE '%{{name}}%'
{{/if}};

-- name: update-user
UPDATE users
SET first_name = ?, last_name = ?
WHERE id = ?;

-- name: delete-user
DELETE FROM users
WHERE id = ?;
