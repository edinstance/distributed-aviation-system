package search.consumer;

import search.config.SearchLogger;
import search.model.FlightDocument;
import search.service.FlightSearchService;
import java.time.Instant;
import lombok.RequiredArgsConstructor;
import org.apache.avro.Schema;
import org.apache.avro.generic.GenericRecord;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

/**
 * This is a kafka consumer for created flights.
 */
@Component
@RequiredArgsConstructor
public class FlightEventConsumer {

  private final FlightSearchService flightSearchService;

  /**
   * This handles the event from kafka.
   *
   * @param flightCreated the flight that was created.
   */
  @KafkaListener(topics = "flights", groupId = "search-service")
  public void handleFlightCreated(GenericRecord flightCreated) {
    try {
      SearchLogger.info("Received FlightCreated event with schema: {} and fields: {}",
              flightCreated.getSchema().getName(),
              flightCreated.getSchema().getFields().stream().map(Schema.Field::name).toList());
      SearchLogger.info("Received FlightCreated event data: {}", flightCreated.toString());

      Object flightIdObj = flightCreated.get("flightId");
      if (flightIdObj == null) {
        SearchLogger.error("FlightId field is null in the received event");
        return;
      }
      String flightId = flightIdObj.toString();
      SearchLogger.info("Received FlightCreated event for flight: {}", flightId);

      FlightDocument flightDocument = FlightDocument.builder()
              .flightId(flightId)
              .number(flightCreated.get("number").toString())
              .origin(flightCreated.get("origin").toString())
              .destination(flightCreated.get("destination").toString())
              .route(flightCreated.get("origin") + " -> " + flightCreated.get("destination"))
              .departureTime(flightCreated.get("departureTime").toString())
              .arrivalTime(flightCreated.get("arrivalTime").toString())
              .airline(flightCreated.get("airline").toString())
              .status(flightCreated.get("status").toString())
              .indexedAt(Instant.now().toString())
              .build();

      flightSearchService.indexFlight(flightDocument);
      SearchLogger.info("Successfully indexed flight: {}", flightId);

    } catch (Exception e) {
      SearchLogger.error("Failed to process FlightCreated event: {}", e.getMessage(), e);
    }
  }
}