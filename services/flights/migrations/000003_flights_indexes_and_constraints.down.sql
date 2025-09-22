ALTER TABLE flights
    DROP CONSTRAINT IF EXISTS chk_times,
    DROP CONSTRAINT IF EXISTS chk_airports;

DROP INDEX IF EXISTS idx_flights_number;
DROP INDEX IF EXISTS idx_flights_route_departure;