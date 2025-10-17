package aviation.aircraft.aircraft.entityFetchers;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.entityfetchers.AircraftEntityFetcher;
import aviation.aircraft.aircraft.services.AircraftService;
import aviation.aircraft.aircraft.types.AircraftStatus;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
public class SetupEntityFetchersTests {

  @Mock
  public AircraftService aircraftService;

  @InjectMocks
  public AircraftEntityFetcher aircraftEntityFetcher;

  public AircraftEntity aircraftEntity;
  public UUID aircraftId = UUID.randomUUID();

  @BeforeEach
  public void setUp() {
    aircraftEntity = new AircraftEntity(
            aircraftId,
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE,
            UUID.randomUUID(),
            UUID.randomUUID(),
            UUID.randomUUID()
    );
  }
}
