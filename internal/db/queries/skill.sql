-- name: CreateSkill :one
INSERT INTO skills (name, description, category, image, hex_theme_color, cv_profile_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSkill :one
SELECT *
FROM skills
WHERE id = $1;

-- name: ListSkills :many
SELECT *
FROM skills
GROUP BY category, id
LIMIT $1 OFFSET $2;
