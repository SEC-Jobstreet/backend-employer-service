INSERT INTO employer_profile (
    email,
    first_name,
    last_name,
    phone,
    address
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: AssociateEmployerWithEnterprise :exec
INSERT INTO employer_enterprise (
    employer_id,
    enterprise_id
) VALUES (
    $1, $2
);

-- name: UpdateEmailConfirmed :exec
UPDATE employer_profile
SET email_confirmed = true
WHERE email = $1
RETURNING employer_profile.email;

-- name: GetEmployerOne :one
SELECT * FROM employer_profile
WHERE id = $1 LIMIT 1;

-- name: GetEmployerProfileByEnterpriseID :one
SELECT employer_profile.* 
FROM employer_profile
JOIN employer_enterprise ON employer_profile.id = employer_enterprise.employer_id
WHERE employer_enterprise.enterprise_id = $1;

-- name: ListEmployerProfiles :many
SELECT * FROM employer_profile
ORDER BY id
LIMIT $1 OFFSET $2;  -- Optional pagination

-- name: CreateEmployerProfile :one
INSERT INTO employer_profile (
  email,
  first_name,
  last_name,
  phone,
  address
)
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateEmployerProfile :one
UPDATE employer_profile
SET
  email = $1,
  first_name = $2,
  last_name = $3,
  phone = $4,
  address = $5,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = $6
RETURNING *;
