-- name: GetFoundPetPercentage :one
SELECT CASE
         WHEN COUNT(*) = 0 THEN CAST(0 AS FLOAT)
         ELSE CAST(COUNT(status = 'found')  / COUNT(*) AS FLOAT)
       END
FROM missing_pets;

-- name: GetLostPetsCount :one
SELECT COUNT(status='missing') FROM missing_pets;

-- name: GetActiveSpotter :one 
SELECT COUNT(role='user') from users;

-- name: GetAvgDailyEvent :one 
SELECT  from users;
