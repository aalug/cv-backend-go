-- name: CreateProject :one
INSERT INTO projects (title,
                      short_description,
                      description,
                      image,
                      technologies_used,
                      hex_theme_color,
                      project_url,
                      cv_profile_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetProject :one
SELECT *
FROM projects
WHERE id = $1;

-- name: ListProjects :many
SELECT *
FROM projects
ORDER BY technologies_used
LIMIT $1 OFFSET $2;