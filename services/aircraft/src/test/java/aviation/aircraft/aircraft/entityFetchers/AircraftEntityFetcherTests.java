package aviation.aircraft.aircraft.entityFetchers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import org.junit.jupiter.api.Test;

public class AircraftEntityFetcherTests extends SetupEntityFetchersTests {

  @Test
  public void testFetchAircraftSuccess() {
    when(aircraftService.getAircraftById(aircraftId)).thenReturn(Optional.of(aircraftEntity));

    Map<String, Object> values = new HashMap<>();
    values.put("id", aircraftId.toString());

    Optional<AircraftEntity> result = aircraftEntityFetcher.fetchAircraft(values);

    assertTrue(result.isPresent());
    assertEquals(aircraftId, result.get().getId());
  }

  @Test
  public void testFetchAircraftNullId() {
    Map<String, Object> values = new HashMap<>();
    values.put("id", null);

    Optional<AircraftEntity> result = aircraftEntityFetcher.fetchAircraft(values);

    assertFalse(result.isPresent());
  }
}
