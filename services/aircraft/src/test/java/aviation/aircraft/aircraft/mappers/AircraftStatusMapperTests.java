package aviation.aircraft.aircraft.mappers;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

import aviation.aircraft.aircraft.mapper.AircraftStatusMapper;
import aviation.aircraft.aircraft.types.AircraftStatus;
import org.junit.jupiter.api.Test;

class AircraftStatusMapperTests {

  @Test
  void testToProtoNull() {
    assertEquals(AircraftStatusMapper.toProto(null), generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_UNSPECIFIED);
  }

  @Test
  void testFromProtoNull() {
    assertNull(AircraftStatusMapper.fromProto(null));
  }

  @Test
  void testToProtoWithDomainValues() {
    assertEquals(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_AVAILABLE,
            AircraftStatusMapper.toProto(AircraftStatus.AVAILABLE));
    assertEquals(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_IN_SERVICE,
            AircraftStatusMapper.toProto(AircraftStatus.IN_SERVICE));
    assertEquals(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_MAINTENANCE,
            AircraftStatusMapper.toProto(AircraftStatus.MAINTENANCE));
    assertEquals(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_GROUNDED,
            AircraftStatusMapper.toProto(AircraftStatus.GROUNDED));
  }

  @Test
  void testFromProtoWithProtoValues() {
    assertEquals(AircraftStatus.AVAILABLE,
            AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_AVAILABLE));
    assertEquals(AircraftStatus.IN_SERVICE,
            AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_IN_SERVICE));
    assertEquals(AircraftStatus.MAINTENANCE,
            AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_MAINTENANCE));
    assertEquals(AircraftStatus.GROUNDED,
            AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_GROUNDED));
  }

  @Test
  void testFromProtoWithSpecialCases() {
    assertNull(AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.AIRCRAFT_STATUS_UNSPECIFIED));
    assertNull(AircraftStatusMapper.fromProto(generated.aircraft.v1.AircraftStatus.UNRECOGNIZED));
  }
}