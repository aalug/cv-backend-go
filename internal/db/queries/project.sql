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

-- name: ListProjectsBySkillName :many
SELECT p.id,
       p.title,
       p.short_description,
       p.description,
       p.image,
       p.hex_theme_color,
       p.project_url,
       p.significance
FROM projects p
         JOIN project_skills ps ON p.id = ps.project_id
         JOIN skills s ON ps.skill_id = s.id
WHERE s.name = sqlc.arg(skill_name)::text
  AND p.cv_profile_id = $1
ORDER BY significance
LIMIT $2 OFFSET $3;
