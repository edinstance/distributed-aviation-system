package aviation.aircraft.aircraft.mutations;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.CreateAircraftMapper;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsMutation;
import com.netflix.graphql.dgs.InputArgument;

/**
 * The aircraft mutations.
 */
@DgsComponent
public class AircraftMutations {

  private final AircraftService aircraftService;

  /**
   * Mutations constructor with an aircraft service.
   *
   * @param aircraftService the service for the mutations to use.
   */
  public AircraftMutations(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  /**
   * A mutation to create an aircraft.
   *
   * @param input the input of the aircraft to be created.
   * @return the created aircraft.
   */
  @DgsMutation
  public AircraftEntity createAircraft(@InputArgument CreateAircraftInput input) {
    return aircraftService.createAircraft(CreateAircraftMapper.toEntity(input));
  }
}
