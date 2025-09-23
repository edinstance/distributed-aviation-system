package aviation.aircraft.aircraft.mutations;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.CreateAircraftMapper;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsMutation;
import com.netflix.graphql.dgs.InputArgument;

@DgsComponent
public class AircraftMutations {

  private final AircraftService aircraftService;

  public AircraftMutations(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  @DgsMutation
  public AircraftEntity createAircraft(@InputArgument CreateAircraftInput input) {
    return aircraftService.createAircraft(CreateAircraftMapper.toEntity(input));
  }
}
