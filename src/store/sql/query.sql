-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  id = ?
LIMIT
  1;

-- name: GetUserByUsernamePassword :one
SELECT
  *
FROM
  users
WHERE
  username = ?
  AND password = ?;

-- name: ListUsers :many
SELECT
  *
FROM
  users
ORDER BY
  username;

-- name: CreateUser :one
INSERT INTO
  users (username, email, phone, name, password)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
  username = ?,
  email = ?,
  phone = ?,
  name = ?,
  password = ?
WHERE
  id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = ?;

-- name: ListChats :many
SELECT
  *
FROM
  chats
ORDER BY
  created_at;

-- name: CreateChat :one
INSERT INTO
  chats (user_id, message)
VALUES
  (?, ?) RETURNING *;

-- name: UpdateChat :exec
UPDATE chats
SET
  message = ?
WHERE
  id = ?;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE
  id = ?;