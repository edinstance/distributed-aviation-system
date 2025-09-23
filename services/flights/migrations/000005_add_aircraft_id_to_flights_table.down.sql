ALTER TABLE flights
    DROP COLUMN aircraft_id;

DROP INDEX IF EXISTS idx_aircraft_id;