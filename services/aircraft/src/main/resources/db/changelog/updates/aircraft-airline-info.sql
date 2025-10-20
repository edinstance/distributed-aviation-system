ALTER TABLE aircraft
    ADD COLUMN airline VARCHAR(100) DEFAULT 'System';

CREATE INDEX idx_aircraft_airline ON aircraft (airline);