package aviation.aircraft.aircraft.queries;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import org.junit.jupiter.api.Test;

public class GetAircraftByIdTests extends SetupQueriesTests {

  @Test
  public void testGetAircraftById() {
    when(aircraftService.getAircraftById(aircraftId)).thenReturn(aircraftEntity);

    AircraftEntity result = aircraftQueries.getAircraftById(aircraftId.toString());

    assertNotNull(result);
  }
}
