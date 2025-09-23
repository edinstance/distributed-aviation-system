package aviation.aircraft.aircraft.mapper;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;

public class CreateAircraftMapper {

  public static AircraftEntity toEntity(CreateAircraftInput createAircraftInput) {

    return new AircraftEntity(
            null,
            createAircraftInput.getRegistration(),
            createAircraftInput.getManufacturer(),
            createAircraftInput.getModel(),
            createAircraftInput.getYearOfManufacture(),
            createAircraftInput.getCapacity(),
            createAircraftInput.getStatus()
    );

  }
}
