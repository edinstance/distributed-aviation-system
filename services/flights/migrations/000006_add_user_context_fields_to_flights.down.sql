ALTER TABLE flights
    DROP COLUMN IF EXISTS organization_id,
    DROP COLUMN IF EXISTS last_updated_by,
    DROP COLUMN IF EXISTS created_by;