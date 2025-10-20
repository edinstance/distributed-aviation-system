package aviation.search.service;

import aviation.search.config.SearchLogger;
import aviation.search.helpers.SearchHelpers;
import aviation.search.model.FlightDocument;
import java.io.IOException;
import java.util.List;
import lombok.RequiredArgsConstructor;
import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.FieldValue;
import org.opensearch.client.opensearch._types.query_dsl.Query;
import org.opensearch.client.opensearch.core.IndexRequest;
import org.springframework.stereotype.Service;

/**
 * This is the search service for flights.
 */
@Service
@RequiredArgsConstructor
public class FlightSearchService {

  private final OpenSearchClient openSearchClient;
  private final FlightIndexService flightIndexService;

  /**
   * This function indexes a flight in opensearch.
   *
   * @param flightDocument the flight to index.
   * @throws IOException if an IO exception occurs while indexing.
   */
  public void indexFlight(FlightDocument flightDocument) throws IOException {
    IndexRequest<FlightDocument> request = IndexRequest.of(i -> i
            .index(flightIndexService.getFlightsIndex())
            .id(flightDocument.getFlightId())
            .document(flightDocument)
    );

    var response = openSearchClient.index(request);
    SearchLogger.info("Indexed flight {} with result: {}",
            flightDocument.getFlightId(), response.result());
  }

  /**
   * A function to search for flights.
   *
   * @param searchTerm the term to search for.
   * @return the flights related to the term.
   * @throws IOException errors from opensearch.
   */
  public List<FlightDocument> searchFlights(String searchTerm) throws IOException {
    Query query = Query.of(q -> q
            .multiMatch(mm -> mm
                    .query(searchTerm)
                    .fields("number", "origin", "destination", "airline")
            )
    );

    return SearchHelpers.executeFlightSearch(query, flightIndexService, openSearchClient);
  }

  /**
   * This function finds flights by their route.
   *
   * @param origin the flight's origin.
   * @param destination the flight's destination.
   * @return flights that match the route.
   * @throws IOException an error while using opensearch.
   */
  public List<FlightDocument> searchFlightsByRoute(
          String origin, String destination) throws IOException {

    String route = origin + " -> " + destination;

    Query query = Query.of(q -> q
            .term(t -> t
                    .field("route")
                    .value(FieldValue.of(route))
            )
    );

    return SearchHelpers.executeFlightSearch(query, flightIndexService, openSearchClient);
  }

  /**
   * This function finds flights by their airline.
   *
   * @param airline the airline to search for.
   * @return the flights from the airline.
   * @throws IOException an error while using opensearch.
   */
  public List<FlightDocument> searchFlightsByAirline(String airline) throws IOException {
    Query query = Query.of(q -> q
            .term(t -> t.field("airline").value(FieldValue.of(airline)))
    );

    return SearchHelpers.executeFlightSearch(query, flightIndexService, openSearchClient);
  }
}