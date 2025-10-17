package aviation.aircraft.user.context;

import java.util.List;
import java.util.UUID;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

/**
 * A context for the user.
 */
@Builder
@Getter
@Setter
public class UserContext {

  /**
   * The user's id.
   */
  private UUID userId;

  /**
   * The user's organization id.
   */
  private UUID orgId;

  /**
   * The user's roles.
   */
  private List<String> roles;

  /**
   * Default constructor.
   */
  public UserContext() {
  }

  /**
   * Constructor with information.
   *
   * @param userId the user's id.
   * @param orgId  the user's organization id.
   * @param roles  the user's roles.
   */
  public UserContext(UUID userId, UUID orgId, List<String> roles) {
    this.userId = userId;
    this.orgId = orgId;
    this.roles = roles;
  }
}