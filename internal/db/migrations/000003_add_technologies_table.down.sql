ALTER TABLE projects
    ADD COLUMN technologies_used VARCHAR(255)[];

ALTER TABLE projects
    DROP COLUMN IF EXISTS significance;

DROP TABLE IF EXISTS project_technologies;

DROP TABLE IF EXISTS technologies;


