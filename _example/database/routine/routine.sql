-- Auto-generated at Thu, 11 Apr 2019 16:15:43 CEST

-- name: select-all-users
SELECT * FROM users;

-- name: select-user-by-pk
SELECT * FROM users
WHERE id = :id;

-- name: insert-user
INSERT INTO users (id, first_name, last_name)
VALUES (:id, :first_name, :last_name);

-- name: update-user-by-pk
UPDATE users
SET first_name = :first_name, last_name = :last_name
WHERE id = :id;

-- name: delete-user-by-pk
DELETE FROM users
WHERE id = :id;

