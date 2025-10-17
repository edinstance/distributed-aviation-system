package aviation.aircraft.user.context;

import com.netflix.graphql.dgs.context.DgsCustomContextBuilderWithRequest;
import java.util.Map;
import java.util.UUID;
import org.jetbrains.annotations.Nullable;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Component;
import org.springframework.web.context.request.WebRequest;

/**
 * Builds a user context from the request headers.
 */
@Component
public class UserContextBuilder implements DgsCustomContextBuilderWithRequest<UserContext> {

  @Override
  public UserContext build(@Nullable Map<String, ?> extensions,
                           @Nullable HttpHeaders headers,
                           @Nullable WebRequest webRequest) {

    if (headers == null) {
      return UserContext.builder().build();
    }

    UUID userId = safeUuid(headers.getFirst("x-user-sub"));

    UUID orgId = safeUuid(headers.getFirst("x-org-id"));

    String roles = headers.getFirst("x-user-roles");

    return UserContext.builder()
            .userId(userId)
            .orgId(orgId)
            .roles(roles)
            .build();
  }

  private UUID safeUuid(String value) {
    if (value == null || value.isBlank()) {
      return null;
    }
    try {
      return UUID.fromString(value);
    } catch (IllegalArgumentException ex) {
      return null;
    }
  }
}