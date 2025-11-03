package search.exceptions;

import graphql.GraphQLError;
import graphql.GraphqlErrorBuilder;
import graphql.execution.DataFetcherExceptionHandler;
import graphql.execution.DataFetcherExceptionHandlerParameters;
import graphql.execution.DataFetcherExceptionHandlerResult;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import org.springframework.stereotype.Component;

/**
 * A global graphql exception handler that controls how exceptions are sent to the client.
 */
@Component
public class GlobalGraphQlExceptionHandler implements DataFetcherExceptionHandler {

  @Override
  public CompletableFuture<DataFetcherExceptionHandlerResult> handleException(
          DataFetcherExceptionHandlerParameters handlerParameters) {

    Throwable ex = handlerParameters.getException();
    GraphQLError error;

    if (ex instanceof DomainException domainEx) {
      error = GraphqlErrorBuilder.newError()
              .message(domainEx.getMessage())
              .path(handlerParameters.getPath())
              .location(handlerParameters.getSourceLocation())
              .extensions(Map.of(
                      "code", domainEx.getCode(),
                      "category", domainEx.getCategory().name()
              ))
              .build();
    } else {
      String errorId = UUID.randomUUID().toString();
      error = GraphqlErrorBuilder.newError()
              .message("Internal server error")
              .path(handlerParameters.getPath())
              .location(handlerParameters.getSourceLocation())
              .extensions(Map.of(
                      "code", "INTERNAL_ERROR",
                      "category", "INTERNAL",
                      "errorId", errorId
              ))
              .build();
    }

    return CompletableFuture.completedFuture(
            DataFetcherExceptionHandlerResult.newResult().error(error).build()
    );
  }
}