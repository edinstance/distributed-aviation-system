package aviation.aircraft.aircraft.dto;

import aviation.aircraft.aircraft.types.AircraftStatus;
import lombok.Getter;
import lombok.Setter;

/**
 * This is a dto for the input for the create aircraft mutation.
 */
@Getter
@Setter
public class CreateAircraftInput {
  private String registration;
  private String manufacturer;
  private String model;
  private int yearOfManufacture;
  private int capacity;
  private AircraftStatus status;

  /**
   * Default constructor.
   */
  public CreateAircraftInput() {
  }

  /**
   * Constructor with information.
   *
   * @param registration the aircraft's registration number.
   * @param manufacturer the aircraft's manufacturer.
   * @param model the aircraft's model.
   * @param yearOfManufacture the year the aircraft was manufactured.
   * @param capacity the capacity of the aircraft.
   * @param status the current status of the aircraft.
   */
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
