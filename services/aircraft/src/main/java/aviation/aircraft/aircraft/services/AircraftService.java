package aviation.aircraft.aircraft.services;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.exceptions.DuplicateAircraftException;
import aviation.aircraft.aircraft.repositories.AircraftRepository;
import aviation.aircraft.config.AircraftLogger;
import aviation.aircraft.user.context.UserContext;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.util.Optional;
import java.util.UUID;
import org.springframework.stereotype.Service;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.params.GetExParams;

/**
 * A service for interacting with aircraft's.
 */
@Service
public class AircraftService {

  private final AircraftRepository aircraftRepository;
  private final JedisPool jedisPool;

  private static final int CACHE_TTL_SECONDS = 900;

  private final ObjectMapper objectMapper;

  /**
   * The constructor for the service.
   *
   * @param aircraftRepository the repository for aircraft's.
   * @param jedisPool          a pool for interacting with the cache.
   * @param objectMapper       an object mapper for transforming objects.
   */
  public AircraftService(final AircraftRepository aircraftRepository,
                         final JedisPool jedisPool, ObjectMapper objectMapper) {
    this.aircraftRepository = aircraftRepository;
    this.jedisPool = jedisPool;
    this.objectMapper = objectMapper;
  }

  /**
   * A function for persisting an aircraft.
   *
   * @param aircraft the aircraft to persist.
   *
   * @return the persisted aircraft.
   */
  public AircraftEntity createAircraft(final AircraftEntity aircraft, final UserContext userCtx) {
    UUID aircraftId = UUID.randomUUID();

    aircraft.setId(aircraftId);
    aircraft.setCreatedBy(userCtx.getUserId());
    aircraft.setLastUpdatedBy(userCtx.getUserId());
    aircraft.setOrganizationId(userCtx.getOrgId());
    aircraft.setAirline(userCtx.getOrgName());

    AircraftLogger.info(String.format(
            "User %s (org=%s) is creating aircraft id=%s (registration=%s)",
            userCtx.getUserId(),
            userCtx.getOrgId(),
            aircraftId,
            aircraft.getRegistration()
    ));

    aircraftRepository.findByRegistration(aircraft.getRegistration())
            .ifPresent(foundAircraft -> {
              AircraftLogger.warn("Aircraft with registration="
                      + aircraft.getRegistration() + " already exists");
              throw new DuplicateAircraftException(foundAircraft.getRegistration());
            });

    AircraftEntity savedAircraft = aircraftRepository.save(aircraft);

    try (Jedis jedis = jedisPool.getResource()) {
      jedis.setex("aircraft:" + savedAircraft.getId().toString(),
              CACHE_TTL_SECONDS, objectMapper.writeValueAsString(savedAircraft));

    } catch (Exception e) {
      AircraftLogger.error("Error while saving aircraft to cache", e);
    }

    AircraftLogger.info(String.format(
            "Aircraft created id=%s by user=%s org=%s",
            savedAircraft.getId(),
            userCtx.getUserId(),
            userCtx.getOrgId()
    ));
    return savedAircraft;
  }

  /**
   * A function for getting an aircraft by its id.
   *
   * @param id the id of the aircraft.
   *
   * @return the found aircraft.
   */
  public Optional<AircraftEntity> getAircraftById(final UUID id) {
    String key = "aircraft:" + id;

    try (Jedis jedis = jedisPool.getResource()) {
      String json = jedis.getEx(key, new GetExParams().ex(CACHE_TTL_SECONDS));
      if (json != null) {
        try {
          return Optional.of(objectMapper.readValue(json, AircraftEntity.class));
        } catch (IOException e) {
          AircraftLogger.error("Cache data corrupted for key=" + key, e);
        }
      }
    } catch (Exception e) {
      AircraftLogger.error("Cache read error for key=" + key, e);
    }

    AircraftLogger.info("Aircraft with id %s not in cache, falling back on db", id);
    try {
      Optional<AircraftEntity> dbResult = aircraftRepository.findById(id);

      dbResult.ifPresent(foundAircraft -> {
        try (Jedis jedis = jedisPool.getResource()) {
          jedis.setex(key, CACHE_TTL_SECONDS,
                  objectMapper.writeValueAsString(foundAircraft));
        } catch (Exception cacheEx) {
          AircraftLogger.error("Error writingAircraft to cache for key=" + key, cacheEx);
        }
      });

      AircraftLogger.info("Aircraft with id %s found in db", id);
      return dbResult;
    } catch (Exception e) {
      AircraftLogger.error("Error fetching Aircraft from DB, id=" + id, e);
      return Optional.empty();
    }
  }
}
