ALTER TABLE aircraft
    ADD COLUMN airline VARCHAR(30) DEFAULT 'System';

CREATE INDEX idx_aircraft_airline ON aircraft (airline);