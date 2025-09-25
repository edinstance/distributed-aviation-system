package aviation.aircraft.entityFetchers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.HashMap;
import java.util.Map;
import org.junit.jupiter.api.Test;

public class AircraftEntityFetcherTests extends SetupEntityFetchersTests {

  @Test
  public void testFetchAircraftSuccess() {
    when(aircraftService.getAircraftById(aircraftId)).thenReturn(aircraftEntity);

    Map<String, Object> values = new HashMap<>();
    values.put("id", aircraftId.toString());

    AircraftEntity result = aircraftEntityFetcher.fetchAircraft(values);

    assertNotNull(result);
    assertEquals(aircraftId, result.getId());
  }

  @Test
  public void testFetchAircraftNullId() {
    Map<String, Object> values = new HashMap<>();
    values.put("id", null);

    AircraftEntity result = aircraftEntityFetcher.fetchAircraft(values);

    assertNull(result);
  }
}
