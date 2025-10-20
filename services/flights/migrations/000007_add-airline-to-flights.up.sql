ALTER TABLE flights
    ADD COLUMN airline VARCHAR(100) DEFAULT 'System';

CREATE INDEX IF NOT EXISTS idx_airline ON flights (airline);