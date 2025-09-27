package aviation.aircraft.grpc.handlers;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.AircraftGrpc;
import aviation.aircraft.aircraft.services.AircraftService;
import generated.aircraft.v1.GetAircraftByIdRequest;
import generated.aircraft.v1.GetAircraftByIdResponse;
import io.grpc.stub.StreamObserver;
import java.util.Optional;
import java.util.UUID;
import org.springframework.stereotype.Component;

/**
 * Grpc handler for getting an aircraft by its id.
 */
@Component
public class GetAircraftByIdHandler {

  private final AircraftService aircraftService;

  /**
   * The constructor for the handler.
   *
   * @param aircraftService the service the handler uses.
   */
  public GetAircraftByIdHandler(AircraftService aircraftService) {
    this.aircraftService = aircraftService;
  }

  /**
   * A function that handles getting the aircraft and returning the grpc version of it.
   *
   * @param request the grpc request which contains the aircraft id.
   * @param responseObserver a stream observer which watches for the response.
   */
  public void handle(
          GetAircraftByIdRequest request,
          StreamObserver<GetAircraftByIdResponse> responseObserver) {
    try {
      UUID aircraftId = UUID.fromString(request.getId());
      Optional<AircraftEntity> aircraftEntity = aircraftService.getAircraftById(aircraftId);

      if (aircraftEntity.isEmpty()) {
        responseObserver.onError(
                io.grpc.Status.NOT_FOUND
                        .withDescription("Aircraft not found for id: " + request.getId())
                        .asRuntimeException());
        return;
      }

      GetAircraftByIdResponse response =
              GetAircraftByIdResponse.newBuilder()
                      .setAircraft(AircraftGrpc.toProto(aircraftEntity.get()))
                      .build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();

    } catch (IllegalArgumentException e) {
      responseObserver.onError(
              io.grpc.Status.INVALID_ARGUMENT
                      .withDescription("Invalid UUID format: " + request.getId())
                      .asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
              io.grpc.Status.INTERNAL
                      .withDescription("Internal server error")
                      .asRuntimeException());
    }
  }
}