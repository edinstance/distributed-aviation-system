ALTER TABLE flights
    ADD COLUMN airline VARCHAR(20) DEFAULT 'System';

CREATE INDEX IF NOT EXISTS idx_airline ON flights (airline);