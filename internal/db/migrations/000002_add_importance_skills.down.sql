-- Drop the unique constraint
ALTER TABLE skills
    DROP CONSTRAINT IF EXISTS unique_category_importance;

-- Drop the "importance" column
ALTER TABLE skills
    DROP COLUMN IF EXISTS importance;