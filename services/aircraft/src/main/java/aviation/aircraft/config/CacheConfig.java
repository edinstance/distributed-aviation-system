package aviation.aircraft.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;

/**
 * The cache configuration.
 */
@Configuration
public class CacheConfig {

  @Value("${redis.url}")
  private String redisUrl;

  /**
   * A spring bean for the cache pool.
   *
   * @return the cache pool.
   */
  @Bean
  public JedisPool jedisPool() {
    JedisPoolConfig poolConfig = new JedisPoolConfig();
    poolConfig.setJmxEnabled(false);
    poolConfig.setMaxTotal(20);
    poolConfig.setMaxIdle(10);
    poolConfig.setMinIdle(5);
    poolConfig.setBlockWhenExhausted(true);
    return new JedisPool(poolConfig, redisUrl);
  }
}