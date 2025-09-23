package aviation.aircraft.aircraft.repositories;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.Optional;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;

public interface AircraftRepository extends JpaRepository<AircraftEntity, UUID> {
  Optional<AircraftEntity> findByRegistration(String registration);
}