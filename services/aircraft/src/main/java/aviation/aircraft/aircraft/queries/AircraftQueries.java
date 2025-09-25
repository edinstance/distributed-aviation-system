package aviation.aircraft.aircraft.queries;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsQuery;
import com.netflix.graphql.dgs.InputArgument;
import java.util.UUID;

/**
 * The aircraft graphql queries.
 */
@DgsComponent
public class AircraftQueries {

  /**
   * The aircraft service for the queries.
   */
  private final AircraftService aircraftService;

  /**
   * A constructor for the aircraft queries.
   *
   * @param aircraftService the aircraft service for the queries.
   */
  public AircraftQueries(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  /**
   * A query to get an aircraft by its id.
   *
   * @param input the id input.
   * @return the found aircraft.
   */
  @DgsQuery
  public AircraftEntity getAircraftById(@InputArgument String input) {
    return aircraftService.getAircraftById(UUID.fromString(input));
  }
}
