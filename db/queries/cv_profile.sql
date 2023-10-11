-- name: CreateCvProfile :one
INSERT INTO cv_profiles (name, email, phone, address, linkedin_url, github_url, bio, profile_picture)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetCvProfile :one
SELECT *
FROM cv_profiles
WHERE id = $1;
