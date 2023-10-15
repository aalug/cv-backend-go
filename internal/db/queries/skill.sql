-- name: CreateSkill :one
INSERT INTO skills (name, description, category, importance, image, hex_theme_color, cv_profile_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSkill :one
SELECT *
FROM skills
WHERE id = $1;

-- name: ListSkills :many
SELECT *
FROM skills
WHERE cv_profile_id = $1
GROUP BY category, id
ORDER BY importance
LIMIT $2 OFFSET $3;
