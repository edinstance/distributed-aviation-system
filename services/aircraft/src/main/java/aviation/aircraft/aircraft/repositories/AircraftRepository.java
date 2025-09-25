package aviation.aircraft.aircraft.repositories;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.Optional;
import java.util.UUID;
import org.springframework.data.jpa.repository.JpaRepository;

/**
 * A repository for interacting with aircraft entities.
 */
public interface AircraftRepository extends JpaRepository<AircraftEntity, UUID> {

  /**
   * A function to find an aircraft by its registration.
   *
   * @param registration the registration to query against.
   * @return an aircraft if it is found.
   */
  Optional<AircraftEntity> findByRegistration(String registration);
}