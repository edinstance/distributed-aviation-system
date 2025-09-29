package aviation.aircraft.aircraft.services;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.Optional;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import redis.clients.jedis.params.GetExParams;

public class GetAircraftByIdTests extends SetupServiceTests {

  UUID aircraftId = UUID.randomUUID();

  @BeforeEach
  public void setup() {
    aircraft.setId(aircraftId);
  }

  @Test
  public void testGetAircraftByIdSuccess() throws Exception {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.getEx(anyString(), any(GetExParams.class))).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.of(aircraft));

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertTrue(result.isPresent());
    verify(jedis).setex(startsWith("aircraft:"), anyLong(), anyString());
  }

  @Test
  public void testGetAircraftByIdCacheSuccess() throws Exception {
    when(jedisPool.getResource()).thenReturn(jedis);

    String json = objectMapper.writeValueAsString(aircraft);
    when(jedis.getEx(anyString(), any(GetExParams.class))).thenReturn(json);

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertTrue(result.isPresent());
    assertEquals(aircraft.getRegistration(), result.get().getRegistration());

    verify(jedis).getEx(startsWith("aircraft:"), any(GetExParams.class));
  }

  @Test
  public void testGetAircraftByIdCacheErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.getEx(anyString(), any(GetExParams.class))).thenThrow(new RuntimeException());
    when(jedis.setex(anyString(), anyLong(), anyString())).thenThrow(new RuntimeException());

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.of(aircraft));

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertTrue(result.isPresent());
    verify(jedis).setex(startsWith("aircraft:"), anyLong(), anyString());
  }

  @Test
  public void testGetAircraftByIdNothingFound() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.getEx(anyString(), any(GetExParams.class))).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenReturn(Optional.empty());

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertFalse(result.isPresent());
  }

  @Test
  public void testGetAircraftByIdRepositoryErrors() {
    when(jedisPool.getResource()).thenReturn(jedis);
    when(jedis.getEx(anyString(), any(GetExParams.class))).thenReturn(null);

    when(aircraftRepository.findById(aircraft.getId())).thenThrow(new RuntimeException());

    Optional<AircraftEntity> result = aircraftService.getAircraftById(aircraft.getId());

    assertFalse(result.isPresent());
  }
}