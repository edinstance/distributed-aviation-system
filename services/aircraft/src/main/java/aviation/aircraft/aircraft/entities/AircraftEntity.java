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

/**
 * An entity for storing aircraft details.
 */
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

  @Column(name = "created_by", nullable = false)
  private UUID createdBy;

  @Column(name = "last_updated_by", nullable = false)
  private UUID lastUpdatedBy;

  @Column(name = "organization_id", nullable = false)
  private UUID organizationId;

  @Column(name = "airline", nullable = false)
  private String airline;

  /**
   * The default constructor.
   */
  public AircraftEntity() {
  }

  /**
   * A constructor with everything apart from id.
   *
   * @param registration      the aircraft's registration number.
   * @param manufacturer      the aircraft's manufacturer.
   * @param model             the aircraft's model.
   * @param yearOfManufacture the year the aircraft was manufactured.
   * @param capacity          the capacity of the aircraft.
   * @param status            the current status of the aircraft.
   * @param createdBy         the user who created the aircraft.
   * @param lastUpdatedBy     the user who last updated the aircraft.
   * @param organizationId    the organization id of the aircraft.
   * @param airline           the airline of the aircraft.
   */
  public AircraftEntity(String registration,
                        String manufacturer, String model,
                        int yearOfManufacture, int capacity,
                        AircraftStatus status, UUID createdBy,
                        UUID lastUpdatedBy, UUID organizationId,
                        String airline) {
    this.registration = registration;
    this.manufacturer = manufacturer;
    this.model = model;
    this.yearOfManufacture = yearOfManufacture;
    this.capacity = capacity;
    this.status = status;
    this.createdBy = createdBy;
    this.lastUpdatedBy = lastUpdatedBy;
    this.organizationId = organizationId;
    this.airline = airline;
  }

  /**
   * A constructor with all fields.
   *
   * @param id                the aircraft's id.
   * @param registration      the aircraft's registration number.
   * @param manufacturer      the aircraft's manufacturer.
   * @param model             the aircraft's model.
   * @param yearOfManufacture the year the aircraft was manufactured.
   * @param capacity          the capacity of the aircraft.
   * @param status            the current status of the aircraft.
   * @param createdBy         the user who created the aircraft.
   * @param lastUpdatedBy     the user who last updated the aircraft.
   * @param organizationId    the organization id of the aircraft.
   * @param airline           the airline of the aircraft.
   */
  public AircraftEntity(UUID id, String registration,
                        String manufacturer, String model,
                        int yearOfManufacture, int capacity,
                        AircraftStatus status, UUID createdBy,
                        UUID lastUpdatedBy, UUID organizationId,
                        String airline) {
    this.id = id;
    this.registration = registration;
    this.manufacturer = manufacturer;
    this.model = model;
    this.yearOfManufacture = yearOfManufacture;
    this.capacity = capacity;
    this.status = status;
    this.createdBy = createdBy;
    this.lastUpdatedBy = lastUpdatedBy;
    this.organizationId = organizationId;
    this.airline = airline;
  }

}