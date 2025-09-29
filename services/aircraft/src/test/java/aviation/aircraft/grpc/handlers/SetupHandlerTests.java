package aviation.aircraft.grpc.handlers;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import aviation.aircraft.aircraft.types.AircraftStatus;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
public class SetupHandlerTests {

  @Mock
  public AircraftService aircraftService;

  public UUID aircraftId = UUID.randomUUID();
  public AircraftEntity aircraftEntity;

  @BeforeEach
  public void setUp() {
    aircraftEntity = new AircraftEntity(
            aircraftId,
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );
  }
}
