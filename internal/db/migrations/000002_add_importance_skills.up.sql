-- Add importance column to the skills table
ALTER TABLE skills
    ADD COLUMN importance INTEGER NOT NULL DEFAULT 1;

-- Add a unique constraint for the "name" field
ALTER TABLE skills
    ADD CONSTRAINT unique_name UNIQUE (name);

-- Add a unique constraint on the combination of "category" and "importance"
ALTER TABLE skills
    ADD CONSTRAINT unique_category_importance UNIQUE (category, importance);