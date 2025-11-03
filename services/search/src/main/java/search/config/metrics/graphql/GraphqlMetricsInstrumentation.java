package search.config.metrics.graphql;

import graphql.ExecutionResult;
import graphql.execution.DataFetcherResult;
import graphql.execution.instrumentation.Instrumentation;
import graphql.execution.instrumentation.InstrumentationContext;
import graphql.execution.instrumentation.InstrumentationState;
import graphql.execution.instrumentation.parameters.InstrumentationExecutionParameters;
import graphql.execution.instrumentation.parameters.InstrumentationFieldFetchParameters;
import graphql.schema.DataFetcher;
import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.metrics.DoubleHistogram;
import io.opentelemetry.api.metrics.LongCounter;
import io.opentelemetry.api.metrics.Meter;
import java.util.concurrent.CompletableFuture;
import org.jetbrains.annotations.NotNull;
import org.springframework.stereotype.Component;

/**
 * A GraphQL instrumentation that records metrics for GraphQL requests.
 */
@Component
public class GraphqlMetricsInstrumentation implements Instrumentation {

  private static final String OPERATION = "operation";
  private static final String TYPE = "type";
  private static final String STATUS = "status";
  private static final String FIELD = "field";

  private final LongCounter graphqlRequests;
  private final DoubleHistogram graphqlDuration;
  private final DoubleHistogram fieldDuration;
  private final LongCounter fieldRequests;

  /**
   * A constructor for the instrumentation.
   */
  public GraphqlMetricsInstrumentation() {
    Meter meter = GlobalOpenTelemetry.getMeter("search.graphql");

    this.graphqlRequests = meter.counterBuilder("search_graphql_requests_total")
            .setDescription("Total GraphQL operations")
            .setUnit("1")
            .build();

    this.graphqlDuration = meter.histogramBuilder("search_graphql_request_duration_ms")
            .setDescription("Duration of GraphQL operations")
            .setUnit("ms")
            .build();

    this.fieldDuration = meter.histogramBuilder("search_graphql_field_duration_ms")
            .setDescription("Duration of GraphQL field fetchers")
            .setUnit("ms")
            .build();

    this.fieldRequests = meter.counterBuilder("search_graphql_fields_total")
            .setDescription("Total GraphQL field fetches")
            .setUnit("1")
            .build();
  }

  /**
   * Begins the execution of a GraphQL operation.
   *
   * @param parameters the parameters of the execution.
   * @param state      the state of the execution.
   *
   * @return the instrumentation context.
   */
  @Override
  public InstrumentationContext<ExecutionResult> beginExecution(
          InstrumentationExecutionParameters parameters, InstrumentationState state) {

    long start = System.nanoTime();

    return new InstrumentationContext<>() {
      @Override
      public void onDispatched() {
      }

      @Override
      public void onCompleted(ExecutionResult result, Throwable t) {
        long durationMs = (System.nanoTime() - start) / 1_000_000;

        String opName = GraphqlMetricsHelpers.createSafeOpName(
                parameters.getExecutionInput().getOperationName());
        String opType = GraphqlMetricsHelpers.extractOperationType(
                parameters.getExecutionInput().getQuery(),
                parameters.getExecutionInput().getOperationName());
        String status = (t == null && (result == null || result.getErrors().isEmpty()))
                ? "success" : "failure";

        Attributes attrs = Attributes.builder()
                .put(OPERATION, opName)
                .put(TYPE, opType)
                .put(STATUS, status)
                .build();

        graphqlRequests.add(1, attrs);
        graphqlDuration.record(durationMs, attrs);
      }
    };
  }

  /**
   * Instruments a GraphQL field fetcher.
   *
   * @param dataFetcher the field fetcher.
   * @param parameters  the parameters of the field fetcher.
   * @param state       the state of the field fetcher.
   *
   * @return the instrumented field fetcher.
   */
  @Override
  public @NotNull DataFetcher<?> instrumentDataFetcher(
          DataFetcher<?> dataFetcher,
          InstrumentationFieldFetchParameters parameters,
          InstrumentationState state) {

    if (parameters.isTrivialDataFetcher()) {
      return dataFetcher;
    }

    String fieldTag = GraphqlMetricsHelpers.findDatafetcherTag(parameters);

    return environment -> {
      long start = System.nanoTime();
      try {
        Object result = dataFetcher.get(environment);

        if (result instanceof CompletableFuture<?> future) {
          return future.whenComplete(
                  (r, ex) -> recordFieldMetric(fieldTag, start,
                          (ex == null
                                  && !(r instanceof DataFetcherResult<?> d
                                  && !d.getErrors().isEmpty()))
                                  ? "success" : "failure"));
        }

        if (result instanceof DataFetcherResult<?> dfr) {
          if (dfr.getData() instanceof CompletableFuture<?> cf) {
            CompletableFuture<?> instrumented = cf.whenComplete(
                    (r, ex) -> recordFieldMetric(fieldTag, start,
                            (ex == null && dfr.getErrors().isEmpty()
                                    && !(r instanceof DataFetcherResult<?> d
                                    && !d.getErrors().isEmpty()))
                                    ? "success" : "failure"));

            return DataFetcherResult.newResult()
                    .data(instrumented)
                    .errors(dfr.getErrors())
                    .localContext(dfr.getLocalContext())
                    .extensions(dfr.getExtensions())
                    .build();
          }
          recordFieldMetric(fieldTag, start, dfr.getErrors().isEmpty() ? "success" : "failure");
          return result;
        }

        recordFieldMetric(fieldTag, start, "success");
        return result;

      } catch (Exception ex) {
        recordFieldMetric(fieldTag, start, "failure");
        throw ex;
      }
    };
  }

  /**
   * Records a field metric.
   *
   * @param fieldTag the field tag.
   * @param start    the start time of the field.
   * @param status   the status of the field.
   */
  private void recordFieldMetric(String fieldTag, long start, String status) {
    long durationMs = (System.nanoTime() - start) / 1_000_000;
    String sanitized = GraphqlMetricsHelpers.sanitizeFieldName(fieldTag);

    Attributes attrs = Attributes.builder()
            .put(FIELD, sanitized)
            .put(STATUS, status)
            .build();

    fieldRequests.add(1, attrs);
    fieldDuration.record(durationMs, attrs);
  }
}