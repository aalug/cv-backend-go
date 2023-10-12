-- name: CreateCvEducation :one
INSERT INTO cv_educations (institution, degree, start_date, end_date, cv_profile_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCvEducation :one
SELECT *
FROM cv_educations
WHERE id = $1;

-- name: ListCvEducations :many
SELECT *
FROM cv_educations
LIMIT $1 OFFSET $2;
