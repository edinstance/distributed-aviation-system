package aviation.aircraft.aircraft.services;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.mockito.ArgumentMatchers.anyInt;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.startsWith;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import com.fasterxml.jackson.core.JsonProcessingException;
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

    AircraftEntity result = aircraftService.getAircraftById(aircraft.getId());

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

    AircraftEntity result = aircraftService.getAircraftById(aircraft.getId());

    assertNotNull(result);
    assertEquals(aircraft.getRegistration(), result.getRegistration());

    verify(jedis).get(startsWith("aircraft:"));
    verify(jedis).expire(startsWith("aircraft:"), anyLong());
  }

  @Test
  public void testGetAircraftByIdCacheErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenThrow(new RuntimeException());
    when(jedis.setex(anyString(), anyLong(), anyString())).thenThrow(new RuntimeException());

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.of(aircraft));

    AircraftEntity result = aircraftService.getAircraftById(aircraft.getId());

    assertNotNull(result);

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

    AircraftEntity result = aircraftService.getAircraftById(aircraft.getId());

    assertNull(result);
  }

  @Test
  public void testGetAircraftByIdRepositoryErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.get(anyString())).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenThrow(new RuntimeException());

    AircraftEntity result = aircraftService.getAircraftById(aircraft.getId());

    assertNull(result);
  }
}
