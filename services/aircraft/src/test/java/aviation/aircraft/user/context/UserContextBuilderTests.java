package aviation.aircraft.user.context;

import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpHeaders;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class UserContextBuilderTests {

  private UserContextBuilder builder;

  @BeforeEach
  void setUp() {
    builder = new UserContextBuilder();
  }

  @Test
  void TestSuccess() {

    UUID userId = UUID.randomUUID();
    UUID orgId = UUID.randomUUID();

    HttpHeaders headers = new HttpHeaders();
    headers.add("x-user-sub", userId.toString());
    headers.add("x-org-id", orgId.toString());
    headers.add("x-user-roles", "ADMIN,MANAGER");

    UserContext ctx = builder.build(Map.of(), headers, null);

    assertNotNull(ctx);
    assertEquals(userId, ctx.getUserId());
    assertEquals(orgId, ctx.getOrgId());
    assertEquals("ADMIN,MANAGER", ctx.getRoles());
  }

  @Test
  void TestEmptyHeaders() {
    HttpHeaders headers = new HttpHeaders();

    UserContext ctx = builder.build(Map.of(), headers, null);

    assertNotNull(ctx);
    assertNull(ctx.getUserId());
    assertNull(ctx.getOrgId());
    assertNull(ctx.getRoles());
  }

  @Test
  void TestNullHeaders() {
    UserContext ctx = builder.build(Map.of(), null, null);

    assertNotNull(ctx);
    assertNull(ctx.getUserId());
    assertNull(ctx.getOrgId());
    assertNull(ctx.getRoles());
  }
}