package aviation.search.exceptions;

/**
 * An exception for when a user is unauthorized.
 */
public class UnauthorizedException extends DomainException {

  /**
   * The exception call with a message.
   *
   * @param message the exception message.
   */
  public UnauthorizedException(String message) {
    super(message, "UNAUTHORIZED", ExceptionCategories.UNAUTHORIZED);
  }

  /**
   * The exception call with a default message.
   */
  public UnauthorizedException() {
    this("Unauthorized");
  }
}