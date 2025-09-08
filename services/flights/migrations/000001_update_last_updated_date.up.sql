-- This function is used to update the last_updated_date column of a table
CREATE OR REPLACE FUNCTION update_last_updated_date()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;