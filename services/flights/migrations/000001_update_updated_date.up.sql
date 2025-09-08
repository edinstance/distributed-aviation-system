-- This function is used to update the updated_date column of a table
CREATE OR REPLACE FUNCTION update_updated_date()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;