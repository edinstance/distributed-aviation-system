package aviation.search.service;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

/**
 * A service for indexing flights.
 */
@Service
@RequiredArgsConstructor
public class FlightIndexService {

  private static final String FLIGHTS_INDEX = "flights";

  /**
   * A function to return the flight index.
   *
   * @return the flight index.
   */
  public String getFlightsIndex() {
    return FLIGHTS_INDEX;
  }
}