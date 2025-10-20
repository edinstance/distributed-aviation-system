package aviation.aircraft.user.directives;

import aviation.aircraft.exceptions.UnauthorizedException;
import aviation.aircraft.user.context.UserContext;
import com.netflix.graphql.dgs.context.DgsContext;
import graphql.schema.DataFetcher;
import graphql.schema.DataFetchingEnvironment;

record AuthenticatedFetcher(DataFetcher<?> delegate, boolean requiresUser,
                                    boolean requiresOrg) implements DataFetcher<Object> {

  @Override
  public Object get(DataFetchingEnvironment dfe) throws Exception {
    UserContext user = DgsContext.getCustomContext(dfe);
    if (user == null) {
      throw new UnauthorizedException("No authentication context found");
    }
    if (requiresUser && user.getUserId() == null) {
      throw new UnauthorizedException("User ID missing");
    }
    if (requiresOrg && user.getOrgId() == null) {
      throw new UnauthorizedException("Organization ID missing");
    }
    if (requiresOrg && user.getOrgName() == null || user.getOrgName().isBlank()) {
      throw new UnauthorizedException("Organization name missing");
    }
    return delegate.get(dfe);
  }
}