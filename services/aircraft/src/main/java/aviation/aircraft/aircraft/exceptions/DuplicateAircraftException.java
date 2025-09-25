package aviation.aircraft.aircraft.exceptions;


import aviation.aircraft.common.exceptions.DomainException;
import aviation.aircraft.common.exceptions.ExceptionCategories;

/**
 * An exception for when there is a duplicate aircraft in the system.
 */
public class DuplicateAircraftException extends DomainException {

  /**
   * The exception call with a registration.
   *
   * @param registration the duplicate aircraft's registration.
   */
  public DuplicateAircraftException(String registration) {
    super("Aircraft with registration " + registration + " already exists.",
            "DUPLICATE_AIRCRAFT", ExceptionCategories.VALIDATION);
  }
}