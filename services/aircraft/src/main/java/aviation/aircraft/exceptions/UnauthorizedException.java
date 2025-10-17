package aviation.aircraft.exceptions;

public class UnauthorizedException extends DomainException {
  public UnauthorizedException(String message) {
    super(message, "UNAUTHORIZED", ExceptionCategories.SECURITY);
  }

  public UnauthorizedException() {
    this("Unauthorized");
  }
}