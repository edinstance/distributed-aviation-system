package aviation.aircraft.aircraft.services;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.startsWith;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.Optional;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class GetAircraftByIdTests extends SetupServiceTests {

  UUID aircraftId = UUID.randomUUID();

  @BeforeEach
  public void setup() {
    aircraft.setId(aircraftId);
  }

  @Test
  public void testGetAircraftByIdSuccess() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.of(aircraft));

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertNotNull(result);

    verify(jedis).setex(
            startsWith("aircraft:"),
            anyLong(),
            anyString()
    );
  }

  @Test
  public void testGetAircraftByIdCacheSuccess() throws Exception {
    when(jedisPool.getResource()).thenReturn(jedis);

    String json = objectMapper.writeValueAsString(aircraft);
    when(jedis.get(anyString())).thenReturn(json);

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertTrue(result.isPresent());
    assertEquals(aircraft.getRegistration(), result.get().getRegistration());

    verify(jedis).get(startsWith("aircraft:"));
    verify(jedis).expire(startsWith("aircraft:"), anyLong());
  }

  @Test
  public void testGetAircraftByIdCacheErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenThrow(new RuntimeException());
    when(jedis.setex(anyString(), anyLong(), anyString())).thenThrow(new RuntimeException());

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.of(aircraft));

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertTrue(result.isPresent());

    verify(jedis).setex(
            startsWith("aircraft:"),
            anyLong(),
            anyString()
    );
  }

  @Test
  public void testGetAircraftByIdNothingFound() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.empty());

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertFalse(result.isPresent());
  }

  @Test
  public void testGetAircraftByIdRepositoryErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenThrow(new RuntimeException());

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertFalse(result.isPresent());
  }
}
