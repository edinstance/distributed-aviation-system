package aviation.search.resolver;

import aviation.search.model.FlightDocument;
import aviation.search.service.FlightSearchService;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.DgsQuery;
import com.netflix.graphql.dgs.InputArgument;
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
   * @return the flights found.
   * @throws IOException if there were any errors.
   */
  @DgsQuery
  public List<FlightDocument> searchFlights(@InputArgument String searchTerm) throws IOException {
    return flightSearchService.searchFlights(searchTerm);
  }

  /**
   * The search flight by route  query.
   *
   * @param origin the flight origin.
   * @param destination the flight destination.
   * @return the flights found.
   * @throws IOException if there were any errors.
   */
  @DgsQuery
  public List<FlightDocument> searchFlightsByRoute(
          @InputArgument String origin,
          @InputArgument String destination
  ) throws IOException {
    return flightSearchService.searchFlightsByRoute(origin, destination);
  }

  /**
   * The search flight by airline query.
   *
   * @param airline the flight's airline.
   * @return the flights found.
   * @throws IOException if there were any errors.
   */
  @DgsQuery
  public List<FlightDocument> searchFlightsByAirline(
          @InputArgument String airline) throws IOException {

    return flightSearchService.searchFlightsByAirline(airline);
  }
}