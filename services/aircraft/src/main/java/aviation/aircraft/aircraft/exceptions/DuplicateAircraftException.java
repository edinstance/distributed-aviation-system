package aviation.aircraft.aircraft.exceptions;


import aviation.aircraft.common.exceptions.DomainException;
import aviation.aircraft.common.exceptions.ExceptionCategories;

public class DuplicateAircraftException extends DomainException {
  public DuplicateAircraftException(String registration) {
    super("Aircraft with registration " + registration + " already exists.",
            "DUPLICATE_AIRCRAFT", ExceptionCategories.VALIDATION);
  }
}