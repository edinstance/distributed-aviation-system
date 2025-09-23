package aviation.aircraft.aircraft.mutations;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import aviation.aircraft.aircraft.types.AircraftStatus;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
public class SetupMutationTests {

  @Mock
  public AircraftService aircraftService;

  @InjectMocks
  public AircraftMutations aircraftMutations;

  public AircraftEntity aircraftEntity;
  public UUID aircraftId = UUID.randomUUID();
  public CreateAircraftInput createAircraftInput;

  @BeforeEach
  public void setUp() {
    createAircraftInput = new CreateAircraftInput(
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );

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
