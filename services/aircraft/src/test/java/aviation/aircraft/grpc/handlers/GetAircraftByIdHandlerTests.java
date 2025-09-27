package aviation.aircraft.grpc.handlers;

import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import java.util.Objects;
import java.util.Optional;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;

public class GetAircraftByIdHandlerTests extends SetupHandlerTests {

  @Mock
  private StreamObserver<generated.aircraft.v1.GetAircraftByIdResponse> responseObserver;

  @InjectMocks
  private GetAircraftByIdHandler getAircraftByIdHandler;

  generated.aircraft.v1.GetAircraftByIdRequest request = generated.aircraft.v1.GetAircraftByIdRequest.newBuilder()
          .setId(String.valueOf(aircraftId))
          .build();

  @Test
  public void testSuccess() {
    when(aircraftService.getAircraftById(aircraftId))
            .thenReturn(Optional.of(aircraftEntity));

    getAircraftByIdHandler.handle(request, responseObserver);

    verify(responseObserver, times(1)).onCompleted();
  }

  @Test
  public void testNotFound() {
    when(aircraftService.getAircraftById(aircraftId)).thenReturn(Optional.empty());

    getAircraftByIdHandler.handle(request, responseObserver);

    verify(responseObserver, times(0)).onCompleted();

    verify(responseObserver, times(1)).onError(argThat(throwable ->
            Status.fromThrowable(throwable).getCode() == Status.Code.NOT_FOUND
                    && Objects.equals(Status.fromThrowable(throwable).getDescription(),
                    "Aircraft not found for id: " + aircraftId)
    ));
  }

  @Test
  public void testIllegalArgument() {
    when(aircraftService.getAircraftById(aircraftId)).thenThrow(IllegalArgumentException.class);

    getAircraftByIdHandler.handle(request, responseObserver);

    verify(responseObserver, times(0)).onCompleted();
    verify(responseObserver, times(1)).onError(argThat(throwable ->
            Status.fromThrowable(throwable).getCode() == Status.Code.INVALID_ARGUMENT
                    && Objects.equals(Status.fromThrowable(throwable).getDescription(),
                    "Invalid UUID format: " + aircraftId)
    ));
  }

  @Test
  public void testInternalError() {
    when(aircraftService.getAircraftById(aircraftId)).thenThrow(RuntimeException.class);

    getAircraftByIdHandler.handle(request, responseObserver);

    verify(responseObserver, times(0)).onCompleted();
    verify(responseObserver, times(1)).onError(argThat(throwable ->
            Status.fromThrowable(throwable).getCode() == Status.Code.INTERNAL
                    && Objects.equals(Status.fromThrowable(throwable).getDescription(),
                    "Internal server error")
    ));
  }
}
