package aviation.aircraft.aircraft.mutations;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mockStatic;
import static org.mockito.Mockito.when;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import aviation.aircraft.user.context.UserContext;
import com.netflix.graphql.dgs.context.DgsContext;
import graphql.GraphQLContext;
import org.junit.jupiter.api.Test;
import org.mockito.MockedStatic;

public class CreateAircraftTests extends SetupMutationTests {

  @Test
  public void createAircraftSuccess() {

    when(aircraftService.createAircraft(any(AircraftEntity.class), any(UserContext.class))).thenReturn(aircraftEntity);

    try (MockedStatic<DgsContext> dgsContextMock = mockStatic(DgsContext.class)) {
      dgsContextMock.when(() -> DgsContext.getCustomContext(dfe))
              .thenReturn(userContext);

      AircraftEntity result = aircraftMutations.createAircraft(createAircraftInput, dfe);

      assertEquals(aircraftEntity, result);
    }
  }
}
