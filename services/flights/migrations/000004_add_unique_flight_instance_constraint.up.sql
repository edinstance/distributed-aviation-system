ALTER TABLE flights
    ADD CONSTRAINT unique_flight_instance UNIQUE (number, departure_time);