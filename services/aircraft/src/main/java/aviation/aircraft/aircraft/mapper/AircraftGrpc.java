package aviation.aircraft.aircraft.mapper;


import aviation.aircraft.aircraft.entities.AircraftEntity;
import generated.aircraft.v1.Aircraft;

/**
 * An aircraft grpc mapper.
 */
public final class AircraftGrpc {

  /**
   * Default constructor.
   */
  private AircraftGrpc() {
  }

  /**
   * A mapper from entity to proto.
   *
   * @param entity the entity to map.
   * @return the proto aircraft.
   */
  public static Aircraft toProto(AircraftEntity entity) {
    if (entity == null) {
      return Aircraft.getDefaultInstance();
    }
    return Aircraft.newBuilder()
            .setId(entity.getId().toString())
            .setRegistration(entity.getRegistration())
            .setManufacturer(entity.getManufacturer())
            .setModel(entity.getModel())
            .setYearOfManufacture(entity.getYearOfManufacture())
            .setCapacity(entity.getCapacity())
            .setStatus(AircraftStatusMapper.toProto(entity.getStatus()))
            .setAirline(entity.getAirline())
            .build();
  }
}