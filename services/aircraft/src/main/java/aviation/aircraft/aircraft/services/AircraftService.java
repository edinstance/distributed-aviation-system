package aviation.aircraft.aircraft.services;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.exceptions.DuplicateAircraftException;
import aviation.aircraft.aircraft.repositories.AircraftRepository;
import aviation.aircraft.common.config.AircraftLogger;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.UUID;
import org.springframework.stereotype.Service;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.params.SetParams;

@Service
public class AircraftService {

  private final AircraftRepository aircraftRepository;
  private final JedisPool jedisPool;

  private final ObjectMapper objectMapper = new ObjectMapper();

  public AircraftService(final AircraftRepository aircraftRepository,
                         final JedisPool jedisPool) {
    this.aircraftRepository = aircraftRepository;
    this.jedisPool = jedisPool;
  }

  public AircraftEntity createAircraft(final AircraftEntity aircraft) {

    aircraft.setId(UUID.randomUUID());

    aircraftRepository.findByRegistration(aircraft.getRegistration())
            .ifPresent(foundAircraft -> {
              throw new DuplicateAircraftException(foundAircraft.getRegistration());
            });

    AircraftEntity savedAircraft = aircraftRepository.save(aircraft);

    try (Jedis jedis = jedisPool.getResource()) {
      jedis.set("aircraft:" + savedAircraft.getId().toString(),
              objectMapper.writeValueAsString(aircraft),
              SetParams.setParams().ex(900));

    } catch (JsonProcessingException e) {
      AircraftLogger.error("Error while saving aircraft to cache", e);
    }

    return savedAircraft;
  }
}
