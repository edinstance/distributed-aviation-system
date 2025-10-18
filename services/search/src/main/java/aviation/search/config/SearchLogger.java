package aviation.search.config;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * A logger for this application.
 */
public class SearchLogger {

  private static final Logger logger = LoggerFactory.getLogger("SearchLogger");

  /**
   * Default constructor.
   */
  private SearchLogger() {
  }

  /**
   * An info level log with just a message.
   *
   * @param message the message to log.
   */
  public static void info(String message) {
    logger.info(message);
  }

  /**
   * An info level log with a message and extra parameters.
   *
   * @param message the message to log.
   * @param params extra params to log.
   */
  public static void info(String message, Object... params) {
    logger.info(message, params);
  }

  /**
   * A warn level log with just a message.
   *
   * @param message the message to log.
   */
  public static void warn(String message) {
    logger.warn(message);
  }

  /**
   * A warn level log with a message and extra parameters.
   *
   * @param message the message to log.
   * @param params extra params to log.
   */
  public static void warn(String message, Object... params) {
    logger.warn(message, params);
  }

  /**
   * An error level log with just a message.
   *
   * @param message the message to log.
   */
  public static void error(String message) {
    logger.error(message);
  }

  /**
   * An error level log with a message and throwable.
   *
   * @param message the message to log.
   * @param throwable the throwable that caused the error.
   */
  public static void error(String message, Throwable throwable) {
    logger.error(message, throwable);
  }

  /**
   * An info level log with a message, throwable and extra parameters.
   *
   * @param message the message to log.
   * @param exception the exception to log.
   * @param params extra params to log.
   */
  public static void error(String message, Exception exception, Object... params) {
    logger.error(message, exception, params);
  }

  /**
   * An exception level log with a message and extra parameters.
   *
   * @param message the message to log.
   * @param params extra params to log.
   */
  public static void error(String message, Object... params) {
    logger.error(message, params);
  }
}