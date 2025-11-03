package search.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * This is the model of the opensearch flights.
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class FlightDocument {
  private String flightId;
  private String number;
  private String origin;
  private String destination;
  private String route;
  private String departureTime;
  private String arrivalTime;
  private String airline;
  private String status;
  private String indexedAt;
}