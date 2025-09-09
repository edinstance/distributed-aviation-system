CREATE TABLE IF NOT EXISTS flights (
    id              UUID PRIMARY KEY NOT NULL,
    number          VARCHAR(20) NOT NULL,
    origin          VARCHAR(3)  NOT NULL,
    destination     VARCHAR(3)  NOT NULL,
    departure_time  TIMESTAMPTZ NOT NULL,
    arrival_time    TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_date_flights
    BEFORE UPDATE ON flights
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_date();