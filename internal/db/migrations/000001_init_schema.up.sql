CREATE TABLE cv_profiles
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    phone           VARCHAR(255) NOT NULL,
    address         VARCHAR(255) NOT NULL,
    linkedin_url    VARCHAR(255),
    github_url      VARCHAR(255) NOT NULL,
    bio             TEXT         NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT (NOW()),
    profile_picture VARCHAR(255) NOT NULL
);

CREATE TABLE cv_educations
(
    id            SERIAL PRIMARY KEY,
    institution   VARCHAR(255)                        NOT NULL,
    degree        VARCHAR(255)                        NOT NULL,
    start_date    DATE                                NOT NULL,
    end_date      DATE                                NOT NULL,
    cv_profile_id INTEGER REFERENCES cv_profiles (id) NOT NULL
);

CREATE TABLE skills
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255)                        NOT NULL,
    description     TEXT                                NOT NULL,
    category        VARCHAR(255)                        NOT NULL,
    image           VARCHAR(255)                        NOT NULL,
    hex_theme_color VARCHAR(255)                        NOT NULL,
    cv_profile_id   INTEGER REFERENCES cv_profiles (id) NOT NULL
);

CREATE TABLE projects
(
    id                SERIAL PRIMARY KEY,
    title             VARCHAR(255)                        NOT NULL,
    short_description VARCHAR(255)                        NOT NULL,
    description       TEXT                                NOT NULL,
    image             VARCHAR(255)                        NOT NULL,
    technologies_used VARCHAR(255)[] NOT NULL,
    hex_theme_color   VARCHAR(255)                        NOT NULL,
    project_url       VARCHAR(255)                        NOT NULL,
    cv_profile_id     INTEGER REFERENCES cv_profiles (id) NOT NULL
);

