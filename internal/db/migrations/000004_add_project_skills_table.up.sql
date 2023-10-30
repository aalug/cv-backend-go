CREATE TABLE project_skills
(
    project_id INTEGER REFERENCES projects (id) NOT NULL,
    skill_id   INTEGER REFERENCES skills (id)   NOT NULL,
    PRIMARY KEY (project_id, skill_id)
);

CREATE INDEX idx_skills_name ON skills (name);