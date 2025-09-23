package aviation.aircraft.aircraft.mapper;

import aviation.aircraft.aircraft.dto.CreateAircraftDto;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import java.util.UUID;

public class CreateAircraftMapper {

  public static AircraftEntity toEntity(CreateAircraftDto createAircraftDto) {

    return new AircraftEntity(
            null,
            createAircraftDto.getRegistration(),
            createAircraftDto.getManufacturer(),
            createAircraftDto.getModel(),
            createAircraftDto.getYearOfManufacture(),
            createAircraftDto.getCapacity(),
            createAircraftDto.getStatus()
    );

  }
}
