package aviation.aircraft.user.context;

import com.netflix.graphql.dgs.context.DgsCustomContextBuilderWithRequest;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;
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
    String orgName = headers.getFirst("x-org-name");

    List<String> roles = parseRoles(headers.getFirst("x-user-roles"));

    return UserContext.builder()
            .userId(userId)
            .orgId(orgId)
            .orgName(orgName)
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

  private List<String> parseRoles(String rolesHeader) {
    if (rolesHeader == null || rolesHeader.isBlank()) {
      return List.of();
    }

    return Arrays.stream(rolesHeader.split(","))
            .map(String::trim)
            .filter(r -> !r.isEmpty())
            .collect(Collectors.toList());
  }
}