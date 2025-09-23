package aviation.aircraft.aircraft.mappers;

import aviation.aircraft.aircraft.dto.CreateAircraftInput;
import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.CreateAircraftMapper;
import aviation.aircraft.aircraft.types.AircraftStatus;
import org.junit.jupiter.api.Test;

public class CreateAircraftMapperTests {

  @Test
  public void testCreateAircraftMapper() {
    CreateAircraftInput createAircraftInput = new CreateAircraftInput(
            "reg",
            "manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );

    AircraftEntity entity = CreateAircraftMapper.toEntity(createAircraftInput);

    assert entity.getId() == null;
    assert entity.getRegistration().equals(createAircraftInput.getRegistration());
    assert entity.getManufacturer().equals(createAircraftInput.getManufacturer());
    assert entity.getModel().equals(createAircraftInput.getModel());
    assert entity.getYearOfManufacture() == createAircraftInput.getYearOfManufacture();
    assert entity.getCapacity() == createAircraftInput.getCapacity();
    assert entity.getStatus().equals(createAircraftInput.getStatus());
  }
}
