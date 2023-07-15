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
SELECT *
FROM missing_pets
WHERE status = 'missing' AND
      ST_Distance(lost_in, ST_SetSRID(ST_MakePoint($1,$2), 4326)) <= $3 LIMIT $4 OFFSET $5;

-- name: SetPetAsFound :one
UPDATE missing_pets SET status='found' WHERE id =$1 AND user_id=$2 
RETURNING *;
