-- name: CreateProjectSkill :one
INSERT INTO project_skills
(project_id,
 skill_id)
VALUES ($1, $2)
RETURNING *;
