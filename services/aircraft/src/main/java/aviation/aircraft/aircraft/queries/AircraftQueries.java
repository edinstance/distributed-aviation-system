package aviation.aircraft.aircraft.queries;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.CreateAircraftMapper;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsQuery;
import com.netflix.graphql.dgs.InputArgument;
import java.util.UUID;

@DgsComponent
public class AircraftQueries {

  private final AircraftService aircraftService;

  public AircraftQueries(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  @DgsQuery
  public AircraftEntity getAircraftById(@InputArgument UUID input) {
    return aircraftService.getAircraftById(input);
  }
}
