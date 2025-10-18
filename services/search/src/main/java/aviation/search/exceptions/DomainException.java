package aviation.search.exceptions;

import lombok.Getter;

/**
 * A custom domain exception for the project.
 */
@Getter
public abstract class DomainException extends RuntimeException {
  private final String code;
  private final ExceptionCategories category;

  /**
   * A constructor with a message, code and category.
   *
   * @param message the exception message.
   * @param code the exception code.
   * @param category the category of exception.
   */
  protected DomainException(String message, String code, ExceptionCategories category) {
    super(message);
    this.code = code;
    this.category = category;
  }

  /**
   * A constructor with a message and code.
   *
   * @param message the exception message.
   * @param code the exception code.
   */
  protected DomainException(String message, String code) {
    super(message);
    this.code = code;
    this.category = ExceptionCategories.INTERNAL;
  }
}