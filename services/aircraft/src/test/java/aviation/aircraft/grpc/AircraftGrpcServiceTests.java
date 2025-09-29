package aviation.aircraft.grpc;

import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.verifyNoMoreInteractions;

import aviation.aircraft.grpc.handlers.GetAircraftByIdHandler;
import generated.aircraft.v1.GetAircraftByIdRequest;
import generated.aircraft.v1.GetAircraftByIdResponse;
import io.grpc.stub.StreamObserver;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class AircraftGrpcServiceTests {

  @Mock
  private GetAircraftByIdHandler getAircraftByIdHandler;

  @Mock
  private StreamObserver<GetAircraftByIdResponse> responseObserver;

  @InjectMocks
  private AircraftGrpcService aircraftGrpcService;

  @Test
  void testGetAircraftById() {
    GetAircraftByIdRequest request = GetAircraftByIdRequest.newBuilder()
            .setId("123e4567-e89b-12d3-a456-426614174000")
            .build();

    aircraftGrpcService.getAircraftById(request, responseObserver);

    verify(getAircraftByIdHandler, times(1))
            .handle(request, responseObserver);
    verifyNoMoreInteractions(getAircraftByIdHandler);
  }
}