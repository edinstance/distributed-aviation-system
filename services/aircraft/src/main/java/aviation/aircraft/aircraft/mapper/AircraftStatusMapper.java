package aviation.aircraft.aircraft.mapper;

import aviation.aircraft.aircraft.types.AircraftStatus;

/**
 * An aircraft status mapper.
 */
public final class AircraftStatusMapper {

  /**
   * Default constructor.
   */
  private AircraftStatusMapper() {
  }

  /**
   * Maps domain status to Protobuf status.
   */
  public static generated.aircraft.v1.AircraftStatus toProto(AircraftStatus status) {
    if (status == null) {
      return generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_UNSPECIFIED;
    }
    return switch (status) {
      case AVAILABLE -> generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_AVAILABLE;
      case IN_SERVICE -> generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_IN_SERVICE;
      case MAINTENANCE -> generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_MAINTENANCE;
      case GROUNDED -> generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_GROUNDED;
    };
  }

  /**
   * Maps Protobuf status to domain status.
   */
  public static AircraftStatus fromProto(generated.aircraft.v1.AircraftStatus protoStatus) {
    if (protoStatus == null) {
      return null;
    }
    return switch (protoStatus) {
      case AIRCRAFT_STATUS_AVAILABLE -> AircraftStatus.AVAILABLE;
      case AIRCRAFT_STATUS_IN_SERVICE -> AircraftStatus.IN_SERVICE;
      case AIRCRAFT_STATUS_MAINTENANCE -> AircraftStatus.MAINTENANCE;
      case AIRCRAFT_STATUS_GROUNDED -> AircraftStatus.GROUNDED;
      case UNRECOGNIZED, AIRCRAFT_STATUS_UNSPECIFIED -> null;
    };
  }
}