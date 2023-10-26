-- name: CreateProject :one
INSERT INTO projects (title,
                      short_description,
                      description,
                      image,
                      hex_theme_color,
                      project_url,
                      cv_profile_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListProjects :many
SELECT id,
       title,
       short_description,
       description,
       image,
       hex_theme_color,
       project_url,
       significance
FROM projects
WHERE cv_profile_id = $1
ORDER BY significance
LIMIT $2 OFFSET $3;
