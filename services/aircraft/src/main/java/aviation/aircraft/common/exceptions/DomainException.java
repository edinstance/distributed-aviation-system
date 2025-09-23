package aviation.aircraft.common.exceptions;

import lombok.Getter;

@Getter
public abstract class DomainException extends RuntimeException {
  private final String code;
  private final ExceptionCategories category;

  protected DomainException(String message, String code, ExceptionCategories category) {
    super(message);
    this.code = code;
    this.category = category;
  }

  protected DomainException(String message, String code) {
    super(message);
    this.code = code;
    this.category = ExceptionCategories.INTERNAL;
  }
}