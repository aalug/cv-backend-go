-- name: CreateTechnology :one
INSERT INTO technologies (name, url, order_field)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateProjectTechnology :one
INSERT INTO project_technologies (project_id, technology_id)
VALUES ($1, $2)
RETURNING *;

-- name: ListTechnologiesForProject :many
SELECT t.id,
       t.name,
       t.url
FROM project_technologies pt
         JOIN technologies t ON pt.technology_id = t.id
WHERE pt.project_id = $1
ORDER BY t.order_field;