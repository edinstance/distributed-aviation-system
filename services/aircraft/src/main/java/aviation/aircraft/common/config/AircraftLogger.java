package aviation.aircraft.common.config;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class AircraftLogger {

  private static final Logger logger = LogManager.getLogger("AircraftLogger");

  private AircraftLogger() {
  }

  public static void info(String message) {
    logger.info(message);
  }

  public static void info(String message, Object... params) {
    logger.info(message, params);
  }

  public static void warn(String message) {
    logger.warn(message);
  }

  public static void warn(String message, Object... params) {
    logger.warn(message, params);
  }

  public static void error(String message) {
    logger.error(message);
  }

  public static void error(String message, Throwable throwable) {
    logger.error(message, throwable);
  }

  public static void error(String message, Exception exception, Object... params) {
    logger.error(message, exception, params);
  }

  public static void error(String message, Object... params) {
    logger.error(message, params);
  }
}