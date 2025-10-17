package aviation.aircraft.aircraft.mutations;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.services.AircraftService;
import aviation.aircraft.aircraft.types.AircraftStatus;
import aviation.aircraft.user.context.UserContext;
import com.netflix.graphql.dgs.DgsDataFetchingEnvironment;
import com.netflix.graphql.dgs.context.DgsContext;
import java.util.List;
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

  @Mock
  public DgsDataFetchingEnvironment dfe;

  @InjectMocks
  public AircraftMutations aircraftMutations;

  public AircraftEntity aircraftEntity;
  public UUID aircraftId = UUID.randomUUID();
  public CreateAircraftInput createAircraftInput;
  public UserContext userContext = new UserContext();

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
            AircraftStatus.AVAILABLE,
            UUID.randomUUID(),
            UUID.randomUUID(),
            UUID.randomUUID()
    );

    userContext = UserContext.builder()
            .userId(UUID.randomUUID())
            .orgId(UUID.randomUUID())
            .roles(List.of("ADMIN"))
            .build();
  }
}
