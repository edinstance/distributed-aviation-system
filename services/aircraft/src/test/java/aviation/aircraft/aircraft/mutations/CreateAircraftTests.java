package aviation.aircraft.aircraft.mutations;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.CreateAircraftMapper;
import org.junit.jupiter.api.Test;

public class CreateAircraftTests extends SetupMutationTests {

  @Test
  public void createAircraftSuccess() {

    when(aircraftService.createAircraft(any(AircraftEntity.class))).thenReturn(aircraftEntity);

    AircraftEntity result = aircraftMutations.createAircraft(createAircraftInput);

    assertEquals(aircraftEntity, result);
  }

}
