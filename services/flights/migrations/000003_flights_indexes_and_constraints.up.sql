
ALTER TABLE flights
    ADD CONSTRAINT chk_times CHECK (arrival_time > departure_time),
    ADD CONSTRAINT chk_airports CHECK (origin <> destination);

CREATE INDEX IF NOT EXISTS idx_flights_number ON flights (number);
CREATE INDEX IF NOT EXISTS idx_flights_route_departure ON flights (origin, destination, departure_time);