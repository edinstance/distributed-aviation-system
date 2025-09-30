package aviation.aircraft.config.metrics.graphql;

import graphql.ExecutionResult;
import graphql.execution.DataFetcherResult;
import graphql.execution.instrumentation.Instrumentation;
import graphql.execution.instrumentation.InstrumentationContext;
import graphql.execution.instrumentation.InstrumentationState;
import graphql.execution.instrumentation.parameters.InstrumentationExecutionParameters;
import graphql.execution.instrumentation.parameters.InstrumentationFieldFetchParameters;
import graphql.schema.DataFetcher;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import java.util.concurrent.CompletableFuture;
import org.jetbrains.annotations.NotNull;
import org.springframework.stereotype.Component;

/**
 * A GraphQL instrumentation that records metrics for GraphQL requests.
 */
@Component
public class GraphqlMetricsInstrumentation implements Instrumentation {

  private final MeterRegistry registry;

  /**
   * Constructor for the instrumentation.
   *
   * @param registry the registry to record metrics in.
   */
  public GraphqlMetricsInstrumentation(MeterRegistry registry) {
    this.registry = registry;
  }

  /**
   * Begins the execution of a GraphQL request.
   *
   * @param parameters the parameters for the execution.
   * @param state      the state of the execution.
   *
   * @return the instrumentation context.
   */
  @Override
  public InstrumentationContext<ExecutionResult> beginExecution(
          InstrumentationExecutionParameters parameters, InstrumentationState state) {

    Timer.Sample requestSample = Timer.start(registry);

    return new InstrumentationContext<>() {
      @Override
      public void onDispatched() {
      }

      @Override
      public void onCompleted(ExecutionResult result, Throwable t) {
        String opName =
                GraphqlMetricsHelpers.createSafeOpName(
                        parameters.getExecutionInput().getOperationName());

        String opType =
                GraphqlMetricsHelpers.extractOperationType(
                        parameters.getExecutionInput().getQuery(),
                        parameters.getExecutionInput().getOperationName());

        String status =
                (t == null && (result == null || result.getErrors().isEmpty()))
                        ? "success"
                        : "failure";

        requestSample.stop(
                Timer.builder("aircraft_graphql_request_duration_seconds")
                        .publishPercentileHistogram(true)
                        .publishPercentiles(0.5, 0.9, 0.95, 0.99)
                        .tags("operation", opName, "type", opType, "status", status)
                        .description("Duration of complete GraphQL operations")
                        .register(registry));

        registry
                .counter(
                        "aircraft_graphql_requests_total",
                        "operation",
                        opName,
                        "type",
                        opType,
                        "status",
                        status)
                .increment();
      }
    };
  }

  /**
   * Instruments a data fetcher.
   *
   * @param dataFetcher the data fetcher to instrument.
   * @param parameters  the parameters for the instrumentation.
   * @param state       the state of the instrumentation.
   *
   * @return the instrumented data fetcher.
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
      Timer.Sample sample = Timer.start(registry);
      try {
        Object result = dataFetcher.get(environment);

        if (result instanceof CompletableFuture<?> future) {
          return future.whenComplete((r, ex) -> {
            boolean hasErrors = ex == null
                    && r instanceof DataFetcherResult<?>
                    && !((DataFetcherResult<?>) r).getErrors().isEmpty();

            String status = (ex == null && !hasErrors) ? "success" : "failure";
            GraphqlMetricsHelpers.stopFieldTimer(registry, sample, fieldTag, status);
          });
        }

        if (result instanceof DataFetcherResult<?> dfr) {
          if (dfr.getData() instanceof CompletableFuture<?> cf) {
            boolean outerHasErrors = !dfr.getErrors().isEmpty();

            CompletableFuture<?> instrumented = cf.whenComplete((r, ex) -> {
              boolean innerHasErrors = ex == null
                      && r instanceof DataFetcherResult<?>
                      && !((DataFetcherResult<?>) r).getErrors().isEmpty();

              String status = (ex == null && !outerHasErrors && !innerHasErrors)
                      ? "success"
                      : "failure";
              GraphqlMetricsHelpers.stopFieldTimer(registry, sample, fieldTag, status);
            });

            return DataFetcherResult.newResult()
                    .data(instrumented)
                    .errors(dfr.getErrors())
                    .localContext(dfr.getLocalContext())
                    .extensions(dfr.getExtensions())
                    .build();
          }

          String status = dfr.getErrors().isEmpty() ? "success" : "failure";
          GraphqlMetricsHelpers.stopFieldTimer(registry, sample, fieldTag, status);
          return result;
        }

        GraphqlMetricsHelpers.stopFieldTimer(registry, sample, fieldTag, "success");
        return result;

      } catch (Exception ex) {
        GraphqlMetricsHelpers.stopFieldTimer(registry, sample, fieldTag, "failure");
        throw ex;
      }
    };
  }
}