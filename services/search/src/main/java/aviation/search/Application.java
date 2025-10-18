package aviation.search;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.jdbc.DataSourceAutoConfiguration;
import org.springframework.boot.autoconfigure.liquibase.LiquibaseAutoConfiguration;

/**
 * This is the main entrypoint to the application.
 */
@SpringBootApplication(exclude = {DataSourceAutoConfiguration.class, LiquibaseAutoConfiguration.class})
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