-- name: CreateApplication :one
INSERT INTO applications (
    candidate_id, 
    job_id,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetApplication :one
SELECT * FROM applications
WHERE id = $1 LIMIT 1;

-- name: ListApplications :many
SELECT * FROM applications
WHERE 
    (candidate_id = @candidate_id OR @candidate_id = 0)
    AND (job_id = @job_id OR @job_id = 0)
    AND (status = @status OR @status = '')
ORDER BY id
LIMIT @l
OFFSET @o;

-- name: UpdateStatusApplication :one
UPDATE applications
SET
    status = sqlc.arg(status),
    message = COALESCE(sqlc.narg(message), message),
    updated_at = sqlc.arg(updated_at)
WHERE
    id = sqlc.arg(id)
RETURNING *;