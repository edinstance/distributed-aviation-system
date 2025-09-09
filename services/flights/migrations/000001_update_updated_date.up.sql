-- This function updates the updated_at column on row modification
CREATE OR REPLACE FUNCTION update_updated_date()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;