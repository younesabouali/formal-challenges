-- name: CreateMissingPet :one
INSERT INTO missing_pets (
  id, createdAt,updatedAt,
  pet_name,description, image_url,
  status,lost_in,lost_at,
  user_id
)
VALUES ($1,$2,$3,
  $4,$5, $6,
  $7,$8,$9,
  $10)
returning *;

-- name: GetMissingPets :many
SELECT * FROM missing_pets LIMIT $1 OFFSET $2;

-- name: SetPetAsFound :one
UPDATE missing_pets SET status='found' WHERE id =$1 
RETURNING *;
