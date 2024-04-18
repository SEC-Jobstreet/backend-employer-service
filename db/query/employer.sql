INSERT INTO employer_profile (
    enterprise_id, 
    email,
    first_name,
    last_name,
    phone,
    address
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetEmployerOne :one
SELECT * FROM employer_profile
WHERE id = $1 LIMIT 1;

-- name: GetEmployerProfileByEnterpriseID :one
SELECT * FROM employer_profile
WHERE enterprise_id = $1
LIMIT 1;

-- name: ListEmployerProfiles :many
SELECT * FROM employer_profile
ORDER BY id
LIMIT $1 OFFSET $2;  -- Optional pagination

-- name: CreateEmployerProfile :one
INSERT INTO employer_profile (
  enterprise_id,
  email,
  first_name,
  last_name,
  phone,
  address
)
VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateEmployerProfile :one
UPDATE employer_profile
SET
  enterprise_id = $1,
  email = $2,
  first_name = $3,
  last_name = $4,
  phone = $5,
  address = $6,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = $7
RETURNING *;
