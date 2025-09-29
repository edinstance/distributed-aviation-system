package aviation.aircraft.aircraft.dto;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import aviation.aircraft.aircraft.types.AircraftStatus;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class CreateAircraftInputTests {

  public CreateAircraftInput createAircraftInput;

  @Test
  public void testCreateAircraftInputDefaultConstructor() {
    createAircraftInput = new CreateAircraftInput();

    assertNotNull(createAircraftInput);
  }

  @Test
  public void testCreateAircraftInputConstructor() {
    createAircraftInput = new CreateAircraftInput(
            "reg",
            "manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );
    assertNotNull(createAircraftInput);

    assertEquals("reg", createAircraftInput.getRegistration());
    assertEquals("manufacturer", createAircraftInput.getManufacturer());
    assertEquals("A380", createAircraftInput.getModel());
    assertEquals(2020, createAircraftInput.getYearOfManufacture());
    assertEquals(50, createAircraftInput.getCapacity());
    assertEquals(AircraftStatus.AVAILABLE, createAircraftInput.getStatus());
  }

  @BeforeEach
  void setUp() {
    createAircraftInput = new CreateAircraftInput(
            "reg",
            "manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );
  }

  @Test
  void testRegistrationMethods() {
    createAircraftInput.setRegistration("New reg");
    assertEquals("New reg", createAircraftInput.getRegistration());
  }

  @Test
  void testManufacturerMethods() {
    createAircraftInput.setManufacturer("New manufacturer");
    assertEquals("New manufacturer", createAircraftInput.getManufacturer());
  }

  @Test
  void testYearOfManufactureMethods() {
    createAircraftInput.setYearOfManufacture(2021);
    assertEquals(2021, createAircraftInput.getYearOfManufacture());
  }

  @Test
  void testCapacityMethods() {
    createAircraftInput.setCapacity(55);
    assertEquals(55, createAircraftInput.getCapacity());
  }

  @Test
  void testStatusMethods() {
    createAircraftInput.setStatus(AircraftStatus.IN_SERVICE);
    assertEquals(AircraftStatus.IN_SERVICE, createAircraftInput.getStatus());
  }
}