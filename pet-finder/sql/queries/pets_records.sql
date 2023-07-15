-- name: RecordEvent :one
INSERT INTO pets_records(
  id, createdAt, updatedAt,
  pet_id,user_id,
  event_name,image_url,description,
  event_location,event_time,area
)
VALUES (
  $1,$2,$3,
  $4,$5,
  $6,$7,$8,
  $9,$10,$11
)
RETURNING *;

-- name: GetMissingPetsEvent :many
SELECT * FROM missing_pets INNER JOIN pets_records ON pets_records.pet_id=missing_pets.id WHERE missing_pets.id=$1;

