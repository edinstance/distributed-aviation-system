package aviation.aircraft.config.cache;

import org.springframework.boot.actuate.health.Health;
import org.springframework.boot.actuate.health.HealthIndicator;
import org.springframework.stereotype.Component;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;

/**
 * A healthcheck for the cache.
 */
@Component("cache")
public class Healthcheck implements HealthIndicator {

  private final JedisPool jedisPool;

  /**
   * Constructor for the healthcheck.
   *
   * @param jedisPool the pool of the cache to be checked.
   */
  public Healthcheck(JedisPool jedisPool) {
    this.jedisPool = jedisPool;
  }

  @Override
  public Health health() {
    try (Jedis jedis = jedisPool.getResource()) {
      String pong = jedis.ping();
      if ("PONG".equalsIgnoreCase(pong)) {
        return Health.up()
                .withDetail("redis", "Available")
                .withDetail("pong", pong)
                .build();
      } else {
        return Health.down()
                .withDetail("redis", "Unexpected response")
                .withDetail("response", pong)
                .build();
      }
    } catch (Exception e) {
      return Health.down(e).withDetail("redis", "Unavailable").build();
    }
  }
}