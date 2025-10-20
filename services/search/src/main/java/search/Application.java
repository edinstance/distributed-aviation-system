package search;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * This is the main entrypoint to the application.
 */
@SpringBootApplication()
public class Application {

  /**
   * This function runs the springboot application.
   *
   * @param args default arguments.
   */
  public static void main(final String[] args) {
    SpringApplication.run(Application.class, args);
  }

}