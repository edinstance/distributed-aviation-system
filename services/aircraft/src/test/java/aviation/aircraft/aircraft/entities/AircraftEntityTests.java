package aviation.aircraft.aircraft.entities;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

import aviation.aircraft.aircraft.types.AircraftStatus;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class AircraftEntityTests {

  public AircraftEntity aircraftEntity;
  public UUID id = UUID.randomUUID();

  @Test
  public void testAircraftEntityDefaultConstructor() {
    aircraftEntity = new AircraftEntity();

    assertNotNull(aircraftEntity);
  }

  @Test
  public void testAircraftEntityConstructor() {
    aircraftEntity = new AircraftEntity(
            id,
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );

    assertNotNull(aircraftEntity);
    assertEquals(id, aircraftEntity.getId());
    assertEquals("Registration", aircraftEntity.getRegistration());
    assertEquals("Manufacturer", aircraftEntity.getManufacturer());
    assertEquals("A380", aircraftEntity.getModel());
    assertEquals(2020, aircraftEntity.getYearOfManufacture());
    assertEquals(50, aircraftEntity.getCapacity());
    assertEquals(AircraftStatus.AVAILABLE, aircraftEntity.getStatus());
  }

  @Test
  public void testAircraftEntityConstructorNoId() {
    aircraftEntity = new AircraftEntity(
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );

    assertNotNull(aircraftEntity);
    assertNull(aircraftEntity.getId());
    assertEquals("Registration", aircraftEntity.getRegistration());
    assertEquals("Manufacturer", aircraftEntity.getManufacturer());
    assertEquals("A380", aircraftEntity.getModel());
    assertEquals(2020, aircraftEntity.getYearOfManufacture());
    assertEquals(50, aircraftEntity.getCapacity());
    assertEquals(AircraftStatus.AVAILABLE, aircraftEntity.getStatus());
  }

  @BeforeEach
  public void setUp() {
    aircraftEntity = new AircraftEntity(
            id,
            "Registration",
            "Manufacturer",
            "A380",
            2020,
            50,
            AircraftStatus.AVAILABLE
    );
  }

  @Test
  public void testAircraftIdMethods() {
    UUID newId = UUID.randomUUID();
    aircraftEntity.setId(newId);
    assertEquals(newId, aircraftEntity.getId());
  }

  @Test
  public void testAircraftRegistrationMethods() {
    aircraftEntity.setRegistration("New Registration");
    assertEquals("New Registration", aircraftEntity.getRegistration());
  }

  @Test
  public void testAircraftManufacturerMethods() {
    aircraftEntity.setManufacturer("New Manufacturer");
    assertEquals("New Manufacturer", aircraftEntity.getManufacturer());
  }

  @Test
  public void testAircraftModelMethods() {
    aircraftEntity.setModel("New Model");
    assertEquals("New Model", aircraftEntity.getModel());
  }

  @Test
  public void testAircraftYearOfManufactureMethods() {
    aircraftEntity.setYearOfManufacture(2021);
    assertEquals(2021, aircraftEntity.getYearOfManufacture());
  }

  @Test
  public void testAircraftCapacityMethods() {
    aircraftEntity.setCapacity(60);
    assertEquals(60, aircraftEntity.getCapacity());
  }

  @Test
  public void testAircraftStatusMethods() {
    aircraftEntity.setStatus(AircraftStatus.IN_SERVICE);
    assertEquals(AircraftStatus.IN_SERVICE, aircraftEntity.getStatus());
  }
}
