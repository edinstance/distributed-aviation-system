package aviation.aircraft.aircraft.entities;


import aviation.aircraft.aircraft.types.AircraftStatus;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import java.util.UUID;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
@Table(name = "aircraft")
public class AircraftEntity {

  @Id
  @Column(unique = true, nullable = false, name = "id")
  private UUID id;

  @Column(name = "registration", nullable = false, unique = true)
  private String registration;

  @Column(name = "manufacturer", nullable = false)
  private String manufacturer;

  @Column(name = "model", nullable = false)
  private String model;

  @Column(name = "year_manufactured", nullable = false)
  private int yearOfManufacture;

  @Column(name = "seat_capacity", nullable = false)
  private int capacity;

  @Enumerated(EnumType.STRING)
  @Column(name = "status", nullable = false)
  private AircraftStatus status;

  public AircraftEntity() {
  }

  public AircraftEntity( String registration,
                        String manufacturer, String model,
                        int yearOfManufacture, int capacity,
                        AircraftStatus status) {
    this.registration = registration;
    this.manufacturer = manufacturer;
    this.model = model;
    this.yearOfManufacture = yearOfManufacture;
    this.capacity = capacity;
    this.status = status;
  }

  public AircraftEntity(UUID id, String registration,
                        String manufacturer, String model,
                        int yearOfManufacture, int capacity,
                        AircraftStatus status) {
    this.id = id;
    this.registration = registration;
    this.manufacturer = manufacturer;
    this.model = model;
    this.yearOfManufacture = yearOfManufacture;
    this.capacity = capacity;
    this.status = status;
  }

}