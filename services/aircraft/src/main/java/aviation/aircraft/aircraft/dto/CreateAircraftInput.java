package aviation.aircraft.aircraft.dto;

import aviation.aircraft.aircraft.types.AircraftStatus;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class CreateAircraftInput {
  private String registration;
  private String manufacturer;
  private String model;
  private int yearOfManufacture;
  private int capacity;
  private AircraftStatus status;

  public CreateAircraftInput() {
  }

  public CreateAircraftInput(String registration, String manufacturer,
                             String model, int yearOfManufacture,
                             int capacity, AircraftStatus status) {
    this.registration = registration;
    this.manufacturer = manufacturer;
    this.model = model;
    this.yearOfManufacture = yearOfManufacture;
    this.capacity = capacity;
    this.status = status;
  }
}
