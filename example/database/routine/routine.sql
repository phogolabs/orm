-- Auto-generated at Thu, 11 Apr 2019 16:15:43 CEST

-- name: select-all-users
SELECT * FROM users;

-- name: insert-user
INSERT INTO users (id, first_name, last_name)
VALUES (:id, :first_name, :last_name);
