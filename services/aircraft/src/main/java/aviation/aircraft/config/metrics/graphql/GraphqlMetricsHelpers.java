package aviation.aircraft.config.metrics.graphql;

import graphql.execution.instrumentation.parameters.InstrumentationFieldFetchParameters;
import graphql.language.Definition;
import graphql.language.Document;
import graphql.language.OperationDefinition;
import graphql.parser.Parser;
import graphql.schema.GraphQLNonNull;
import graphql.schema.GraphQLObjectType;
import graphql.schema.GraphQLOutputType;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;

/**
 * A helper class for recording metrics for GraphQL requests.
 */
public class GraphqlMetricsHelpers {

  private static final Parser parser = new Parser();

  /**
   * Creates a safe operation name for metrics.
   *
   * @param name the operation name.
   * @return the safe operation name.
   */
  public static String createSafeOpName(String name) {
    if (name == null || name.isBlank()) {
      return "unknown";
    }
    return name.length() > 50 ? name.substring(0, 50) : name;
  }

  /**
   * Extracts the operation type from a query.
   *
   * @param query the query.
   * @param operationName the operation name.
   * @return the operation type.
   */
  public static String extractOperationType(String query, String operationName) {
    if (query == null) {
      return "unknown";
    }

    try {
      Document doc = parser.parseDocument(query);
      for (Definition<?> def : doc.getDefinitions()) {
        if (def instanceof OperationDefinition opDef) {
          if (operationName == null || operationName.equals(opDef.getName())) {
            return opDef.getOperation().name();
          }
        }
      }
    } catch (Exception e) {
      return "unknown";
    }
    return "unknown";
  }

  /**
   * Stops a field timer.
   *
   * @param registry the registry to record metrics in.
   * @param sample the sample to stop the timer with.
   * @param fieldTag the field tag.
   * @param status the status of the request.
   */
  public static void stopFieldTimer(
          MeterRegistry registry, Timer.Sample sample, String fieldTag, String status) {

    String sanitized = sanitizeFieldName(fieldTag);

    sample.stop(
            Timer.builder("aircraft_graphql_field_duration_seconds")
                    .publishPercentileHistogram(true)
                    .publishPercentiles(0.5, 0.9, 0.95, 0.99)
                    .tags("field", sanitized, "status", status)
                    .description("Duration of individual GraphQL field fetchers")
                    .register(registry));

    registry
            .counter("aircraft_graphql_fields_total", "field", sanitized, "status", status)
            .increment();
  }

  /**
   * Finds the datafetcher tag for a field.
   *
   * @param parameters the parameters of the field.
   * @return the tag for the field.
   */
  public static String findDatafetcherTag(InstrumentationFieldFetchParameters parameters) {
    GraphQLOutputType type = parameters.getExecutionStepInfo().getParent().getType();
    GraphQLObjectType parent =
            type instanceof GraphQLNonNull nonNull
                    ? (GraphQLObjectType) nonNull.getWrappedType()
                    : (GraphQLObjectType) type;
    return parent.getName() + "." + parameters.getExecutionStepInfo().getField().getName();
  }

  /**
   * Sanitizes a field name for use in metrics.
   *
   * @param tag the field name.
   * @return the sanitized field name.
   */
  public static String sanitizeFieldName(String tag) {
    if (tag == null || tag.isBlank()) {
      return "unknown";
    }

    return tag.replaceAll("[0-9]+", "N");
  }
}