ALTER TABLE flights
    ADD COLUMN aircraft_id UUID;

CREATE INDEX IF NOT EXISTS idx_aircraft_id ON flights (aircraft_id);