CREATE TABLE cv_profile
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    bio             TEXT         NOT NULL,
    profile_picture VARCHAR(255) NOT NULL
);

CREATE TABLE cv_education
(
    id            SERIAL PRIMARY KEY,
    institution   VARCHAR(255) NOT NULL,
    degree        VARCHAR(255) NOT NULL,
    start_date    DATE         NOT NULL,
    end_date      DATE         NOT NULL,
    cv_profile_id INTEGER REFERENCES cv_profile (id)
);

CREATE TABLE skills
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT         NOT NULL,
    image           VARCHAR(255) NOT NULL,
    hex_theme_color VARCHAR(255) NOT NULL,
    cv_profile_id   INTEGER REFERENCES cv_profile (id)
);

CREATE TABLE projects
(
    id                SERIAL PRIMARY KEY,
    title             VARCHAR(255) NOT NULL,
    short_description VARCHAR(255) NOT NULL,
    description       TEXT         NOT NULL,
    image             VARCHAR(255) NOT NULL,
    technologies_used VARCHAR(255)[] NOT NULL,
    hex_theme_color   VARCHAR(255) NOT NULL,
    project_url       VARCHAR(255) NOT NULL,
    cv_profile_id     INTEGER REFERENCES cv_profile (id)
);

