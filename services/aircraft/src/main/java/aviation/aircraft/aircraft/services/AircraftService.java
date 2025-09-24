package aviation.aircraft.aircraft.services;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.exceptions.DuplicateAircraftException;
import aviation.aircraft.aircraft.repositories.AircraftRepository;
import aviation.aircraft.common.config.AircraftLogger;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.UUID;
import org.springframework.stereotype.Service;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;

@Service
public class AircraftService {

  private final AircraftRepository aircraftRepository;
  private final JedisPool jedisPool;

  private static final int CACHE_TTL_SECONDS = 900;

  private final ObjectMapper objectMapper;

  public AircraftService(final AircraftRepository aircraftRepository,
                         final JedisPool jedisPool, ObjectMapper objectMapper) {
    this.aircraftRepository = aircraftRepository;
    this.jedisPool = jedisPool;
    this.objectMapper = objectMapper;
  }

  public AircraftEntity createAircraft(final AircraftEntity aircraft) {

    aircraft.setId(UUID.randomUUID());

    aircraftRepository.findByRegistration(aircraft.getRegistration())
            .ifPresent(foundAircraft -> {
              throw new DuplicateAircraftException(foundAircraft.getRegistration());
            });

    AircraftEntity savedAircraft = aircraftRepository.save(aircraft);

    try (Jedis jedis = jedisPool.getResource()) {
      jedis.setex("aircraft:" + savedAircraft.getId().toString(),
              CACHE_TTL_SECONDS, objectMapper.writeValueAsString(savedAircraft));

    } catch (Exception e) {
      AircraftLogger.error("Error while saving aircraft to cache", e);
    }

    return savedAircraft;
  }


  public AircraftEntity getAircraftById(final UUID id) {
    String key = "aircraft:" + id;

    try (Jedis jedis = jedisPool.getResource()) {
      String json = jedis.get(key);
      if (json != null) {
        jedis.expire(key, CACHE_TTL_SECONDS);
        return objectMapper.readValue(json, AircraftEntity.class);
      }
    } catch (Exception e) {
      AircraftLogger.error("Error reading value from the cache", e);
    }

    try {
      return aircraftRepository.findById(id)
              .map(foundAircraft -> {
                try (Jedis jedis = jedisPool.getResource()) {
                  jedis.setex(key, CACHE_TTL_SECONDS, objectMapper.writeValueAsString(foundAircraft));
                } catch (Exception cacheEx) {
                  AircraftLogger.error("Error writing value to cache", cacheEx);
                }
                return foundAircraft;
              })
              .orElse(null);

    } catch (Exception e) {
      AircraftLogger.error("Error fetching Aircraft from DB", e);
      return null;
    }
  }
}
