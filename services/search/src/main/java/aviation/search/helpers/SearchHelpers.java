package aviation.search.helpers;

import aviation.search.model.FlightDocument;
import aviation.search.service.FlightIndexService;
import java.io.IOException;
import java.util.List;
import java.util.stream.Collectors;
import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.opensearch._types.query_dsl.Query;
import org.opensearch.client.opensearch.core.SearchRequest;
import org.opensearch.client.opensearch.core.SearchResponse;
import org.opensearch.client.opensearch.core.search.Hit;

/**
 * Helpers for searching.
 */
public class SearchHelpers {

  /**
   * A funtion to help execute flight searches.
   *
   * @param query the query to run.
   * @param flightIndexService the flight index service.
   * @param openSearchClient the opensearch client.
   * @return the flights returned from the query.
   * @throws IOException errors from opensearch.
   */
  public static List<FlightDocument> executeFlightSearch(
          Query query,
          FlightIndexService flightIndexService,
          OpenSearchClient openSearchClient) throws IOException {

    SearchRequest request = SearchRequest.of(s -> s
            .index(flightIndexService.getFlightsIndex())
            .query(query)
            .size(50)
    );

    SearchResponse<FlightDocument> response =
            openSearchClient.search(request, FlightDocument.class);

    return response.hits().hits().stream()
            .map(Hit::source)
            .collect(Collectors.toList());
  }
}