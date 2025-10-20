package aviation.aircraft.aircraft.services;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.repositories.AircraftRepository;
import aviation.aircraft.aircraft.types.AircraftStatus;
import aviation.aircraft.user.context.UserContext;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.netflix.graphql.dgs.DgsDataFetchingEnvironment;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.Spy;
import org.mockito.junit.jupiter.MockitoExtension;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;

@ExtendWith(MockitoExtension.class)
public class SetupServiceTests {

  @Mock
  public AircraftRepository aircraftRepository;

  @Mock
  public JedisPool jedisPool;

  @Mock
  public Jedis jedis;

  @Mock
  public DgsDataFetchingEnvironment dfe;

  @Spy
  public ObjectMapper objectMapper = new ObjectMapper();

  @InjectMocks
  public AircraftService aircraftService;

  public AircraftEntity aircraft;

  public UserContext userContext = new UserContext();

  @BeforeEach
  public void setUp() {
    aircraft = new AircraftEntity(
            null,
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE,
            UUID.randomUUID(),
            UUID.randomUUID(),
            UUID.randomUUID(),
            "British Airways"
    );
  }
}