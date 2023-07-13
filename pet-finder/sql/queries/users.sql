-- name: CreateUser :one
INSERT INTO users (
  id, createdAt,updatedAt,name,email,password,role
)
VALUES ($1,$2,$3,$4,$5, $6,$7)
returning *;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email=$1;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key =$1 ;
