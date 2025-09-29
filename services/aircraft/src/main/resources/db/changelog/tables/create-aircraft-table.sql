
CREATE TABLE aircraft (
    id UUID PRIMARY KEY NOT NULL,
    registration VARCHAR(20) NOT NULL UNIQUE,
    manufacturer VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year_manufactured INTEGER,
    seat_capacity INTEGER,
    status VARCHAR(20) DEFAULT 'AVAILABLE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aircraft_registration ON aircraft(registration);

CREATE TRIGGER set_updated_date_aircraft
    BEFORE UPDATE ON aircraft
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_date();