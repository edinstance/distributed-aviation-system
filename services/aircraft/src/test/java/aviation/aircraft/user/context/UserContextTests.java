package aviation.aircraft.user.context;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.AssertionsKt.assertNotNull;

import java.util.List;
import java.util.UUID;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class UserContextTests {

  private UserContext userContext;

  @Test
  public void testUserContextDefaultConstructor() {
    userContext = new UserContext();

    assertNotNull(userContext);
  }

  @Test
  public void testUserContextConstructor() {
    userContext = new UserContext(UUID.randomUUID(), UUID.randomUUID(), List.of("roles"));

    assertNotNull(userContext);
  }

  @BeforeEach
  public void setup() {
    userContext = new UserContext(UUID.randomUUID(), UUID.randomUUID(), List.of("roles"));
  }

  @Test
  public void testUserIdMethods(){
    UUID newId = UUID.randomUUID();
    userContext.setUserId(newId);
    assertEquals(newId, userContext.getUserId());
  }

  @Test
  public void testOrgIdMethods(){
    UUID newId = UUID.randomUUID();
    userContext.setOrgId(newId);
    assertEquals(newId, userContext.getOrgId());
  }

  @Test
  public void testRolesMethods(){
    userContext.setRoles(List.of("New Roles"));
    assert userContext.getRoles().contains("New Roles");
  }
}
