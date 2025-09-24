package aviation.aircraft.common.config;

import aviation.aircraft.aircraft.entities.AircraftEntity;
import com.netflix.graphql.dgs.DgsComponent;
import com.netflix.graphql.dgs.federation.DefaultDgsFederationResolver;
import jakarta.annotation.PostConstruct;
import java.util.HashMap;
import java.util.Map;
import org.jetbrains.annotations.NotNull;

@DgsComponent
public class FederationResolver extends DefaultDgsFederationResolver {
  private final Map<Class<?>, String> types = new HashMap<>();

  @PostConstruct
  public void init() {
    types.put(AircraftEntity.class, "Aircraft");
  }

  @Override
  public @NotNull Map<Class<?>, String> typeMapping() {
    return types;
  }
}