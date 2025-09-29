package aviation.aircraft.grpc;

import aviation.aircraft.grpc.handlers.GetAircraftByIdHandler;
import generated.aircraft.v1.AircraftServiceGrpc;
import generated.aircraft.v1.GetAircraftByIdRequest;
import generated.aircraft.v1.GetAircraftByIdResponse;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;

/**
 * gRPC service implementation for Aircraft operations.
 */
@GrpcService
public class AircraftGrpcService extends AircraftServiceGrpc.AircraftServiceImplBase {

  private final GetAircraftByIdHandler getAircraftByIdHandler;

  /**
   * Constructor for AircraftGrpcService.
   *
   * @param getAircraftByIdHandler the get aircraft by id gRPC handler
   */
  public AircraftGrpcService(final GetAircraftByIdHandler getAircraftByIdHandler) {
    this.getAircraftByIdHandler = getAircraftByIdHandler;
  }

  @Override
  public void getAircraftById(GetAircraftByIdRequest request,
                              StreamObserver<GetAircraftByIdResponse> responseObserver) {
    getAircraftByIdHandler.handle(request, responseObserver);
  }
}