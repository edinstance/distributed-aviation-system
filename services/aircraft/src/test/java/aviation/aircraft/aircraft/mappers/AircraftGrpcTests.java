package aviation.aircraft.aircraft.mappers;

import static org.junit.jupiter.api.Assertions.assertEquals;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.aircraft.mapper.AircraftGrpc;
import aviation.aircraft.aircraft.types.AircraftStatus;
import generated.aircraft.v1.Aircraft;
import java.util.UUID;
import org.junit.jupiter.api.Test;

class AircraftGrpcTests {

  @Test
  void testNullInputDefaultInstance() {
    Aircraft proto = AircraftGrpc.toProto(null);

    assertEquals(Aircraft.getDefaultInstance(), proto);
  }

  @Test
  void testToProtoMapsAllFields() {

    AircraftEntity entity = new AircraftEntity(
            UUID.randomUUID(),
            "Registration",
            "Manufacturer",
            "Model",
            2020,
            180,
            AircraftStatus.AVAILABLE,
            UUID.randomUUID(),
            UUID.randomUUID(),
            UUID.randomUUID()
    );

    Aircraft proto = AircraftGrpc.toProto(entity);

    assertEquals(entity.getId(), UUID.fromString(proto.getId()));
    assertEquals(entity.getRegistration(), proto.getRegistration());
    assertEquals(entity.getManufacturer(), proto.getManufacturer());
    assertEquals(entity.getModel(), proto.getModel());
    assertEquals(entity.getYearOfManufacture(), proto.getYearOfManufacture());
    assertEquals(entity.getCapacity(), proto.getCapacity());
  }
}