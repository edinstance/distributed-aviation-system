package aviation.aircraft.aircraft.entityfetchers;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsEntityFetcher;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import org.jetbrains.annotations.NotNull;

/**
 * This is an entity fetcher for aircraft which allows federation to query for them.
 */
@DgsComponent
public class AircraftEntityFetcher {

  private final AircraftService aircraftService;

  /**
   * Constructor for the aircraft entity fetcher.
   *
   * @param aircraftService the aircraft service that the entity fetcher can use.
   */
  public AircraftEntityFetcher(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }


  /**
   * The entity fetcher for fetching aircraft's.
   *
   * @param values the params for the search.
   *
   * @return if an entity is found then it is returned if not null is returned.
   */
  @DgsEntityFetcher(name = "Aircraft")
  public Optional<AircraftEntity> fetchAircraft(@NotNull Map<String, Object> values) {
    String id = (String) values.get("id");
    if (id == null) {
      return Optional.empty();
    }

    return aircraftService.getAircraftById(UUID.fromString(id));
  }
}
