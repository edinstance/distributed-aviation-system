package aviation.aircraft.aircraft.mappers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

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

    assertNull(entity.getId());
    assertEquals(createAircraftInput.getRegistration(), entity.getRegistration());
    assertEquals(createAircraftInput.getManufacturer(), entity.getManufacturer());
    assertEquals(createAircraftInput.getModel(), entity.getModel());
    assertEquals(createAircraftInput.getYearOfManufacture(), entity.getYearOfManufacture());
    assertEquals(createAircraftInput.getCapacity(), entity.getCapacity());
    assertEquals(createAircraftInput.getStatus(), entity.getStatus());
  }
}
