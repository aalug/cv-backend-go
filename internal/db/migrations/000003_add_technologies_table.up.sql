CREATE TABLE technologies
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    url         VARCHAR(255) NOT NULL,
    order_field INTEGER      NOT NULL
);

CREATE TABLE project_technologies
(
    project_id    INTEGER REFERENCES projects (id)     NOT NULL,
    technology_id INTEGER REFERENCES technologies (id) NOT NULL,
    PRIMARY KEY (project_id, technology_id)
);

ALTER TABLE projects
    DROP COLUMN technologies_used;

ALTER TABLE projects
    ADD COLUMN significance integer NOT NULL DEFAULT 0;
