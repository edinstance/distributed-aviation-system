package aviation.aircraft.aircraft.mapper;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;

/**
 * A mapper for mapping create aircraft data transfer objects to maps.
 */
public class CreateAircraftMapper {

  /**
   * The function which maps an input to an entity.
   *
   * @param createAircraftInput the create aircraft input that should be mapped.
   *
   * @return the input as an entity.
   */
  public static AircraftEntity toEntity(CreateAircraftInput createAircraftInput) {

    return new AircraftEntity(
            null,
            createAircraftInput.getRegistration(),
            createAircraftInput.getManufacturer(),
            createAircraftInput.getModel(),
            createAircraftInput.getYearOfManufacture(),
            createAircraftInput.getCapacity(),
            createAircraftInput.getStatus(),
            null,
            null,
            null,
            null
    );

  }
}
