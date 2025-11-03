package search.resolver;

import search.model.FlightDocument;
import search.service.FlightSearchService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsQuery;
import com.netflix.graphql.dgs.InputArgument;
import com.netflix.graphql.dgs.exceptions.DgsEntityNotFoundException;
import java.io.IOException;
import java.util.List;
import lombok.RequiredArgsConstructor;

/**
 * The graphql resolvers for searching flights.
 */
@DgsComponent
@RequiredArgsConstructor
public class FlightSearchResolver {

  private final FlightSearchService flightSearchService;

  /**
   * The search flight query.
   *
   * @param searchTerm the term to search for.
   *
   * @return the flights found.
   */
  @DgsQuery
  public List<FlightDocument> searchFlights(@InputArgument String searchTerm) {
    try {
      return flightSearchService.searchFlights(searchTerm);
    } catch (IOException e) {
      throw new DgsEntityNotFoundException("Failed to search flights: " + e.getMessage());
    }
  }

  /**
   * The search flight by route  query.
   *
   * @param origin      the flight origin.
   * @param destination the flight destination.
   *
   * @return the flights found.
   */
  @DgsQuery
  public List<FlightDocument> searchFlightsByRoute(
          @InputArgument String origin,
          @InputArgument String destination
  ) {
    try {
      return flightSearchService.searchFlightsByRoute(origin, destination);
    } catch (IOException e) {
      throw new DgsEntityNotFoundException("Failed to search flights: " + e.getMessage());
    }
  }

  /**
   * The search flight by airline query.
   *
   * @param airline the flight's airline.
   *
   * @return the flights found.
   */
  @DgsQuery
  public List<FlightDocument> searchFlightsByAirline(
          @InputArgument String airline) {
    try {
      return flightSearchService.searchFlightsByAirline(airline);
    } catch (IOException e) {
      throw new DgsEntityNotFoundException("Failed to search flights: " + e.getMessage());
    }
  }
}