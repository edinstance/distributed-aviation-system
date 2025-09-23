package aviation.aircraft.aircraft.services;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.startsWith;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.exceptions.DuplicateAircraftException;
import java.util.Optional;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
public class CreateAircraftTests extends SetupServiceTests {


  @Test
  public void testCreateAircraftSuccess() {
    when(aircraftRepository.findByRegistration(aircraft.getRegistration())).thenReturn(Optional.empty());
    when(aircraftRepository.save(any(AircraftEntity.class))).thenAnswer(inv -> inv.getArgument(0));
    when(jedisPool.getResource()).thenReturn(jedis);

    AircraftEntity result = aircraftService.createAircraft(aircraft);

    assertNotNull(result);
    assertNotNull(result.getId());

    assertEquals(aircraft.getRegistration(), result.getRegistration());

    verify(aircraftRepository).findByRegistration(aircraft.getRegistration());
    verify(aircraftRepository).save(aircraft);

    verify(jedis).set(
            startsWith("aircraft:"),
            anyString(),
            any()
    );
  }

  @Test
  public void testCreateAircraftDuplicate() {
    when(aircraftRepository.findByRegistration(aircraft.getRegistration())).thenReturn(Optional.of(aircraft));

    assertThrows(DuplicateAircraftException.class, () -> aircraftService.createAircraft(aircraft));

    verify(aircraftRepository, times(0)).save(aircraft);
  }
}