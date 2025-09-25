package aviation.aircraft.aircraft.entityFetchers;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsEntityFetcher;
import java.util.Map;
import java.util.UUID;
import org.jetbrains.annotations.NotNull;

@DgsComponent
public class AircraftEntityFetcher {

  private final AircraftService aircraftService;

  public AircraftEntityFetcher(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  @DgsEntityFetcher(name = "Aircraft")
  public AircraftEntity fetchAircraft(@NotNull Map<String, Object> values){
    String id = (String) values.get("id");
    if (id == null) {
      return null;
    }

    return aircraftService.getAircraftById(UUID.fromString(id));
  }
}
