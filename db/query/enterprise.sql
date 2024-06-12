-- name: CreateEnterprise :one
INSERT INTO enterprise (
    id,
    name,
    country, 
    address,
    latitude,
    longitude,
    field,
    size,
    url,
    license,
    employer_id,
    employer_role
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetEnterpriseById :one
SELECT * FROM enterprise
WHERE id = $1 LIMIT 1;

-- name: GetEnterpriseByEmployerId :many
SELECT * FROM enterprise
WHERE employer_id = $1;

-- name: UpdateEnterprise :one
UPDATE enterprise
SET
  name = COALESCE(sqlc.narg(name), name),
  country = COALESCE(sqlc.narg(country), country),
  address = COALESCE(sqlc.narg(address), address),
  field = COALESCE(sqlc.narg(field), field),
  latitude = COALESCE(sqlc.narg(latitude), latitude),
  longitude = COALESCE(sqlc.narg(longitude), longitude),
  size = COALESCE(sqlc.narg(size), size),
  url = COALESCE(sqlc.narg(url), url),
  license = COALESCE(sqlc.narg(license), license),
  employer_role = COALESCE(sqlc.narg(employer_role), employer_role)
WHERE
  id = sqlc.arg(id) and employer_id = sqlc.arg(employer_id)
RETURNING *;