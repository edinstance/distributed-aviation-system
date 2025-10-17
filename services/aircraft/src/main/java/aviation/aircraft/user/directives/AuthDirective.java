package aviation.aircraft.user.directives;

import com.netflix.graphql.dgs.DgsDirective;
import graphql.schema.DataFetcher;
import graphql.schema.GraphQLAppliedDirective;
import graphql.schema.GraphQLAppliedDirectiveArgument;
import graphql.schema.GraphQLFieldDefinition;
import graphql.schema.GraphQLObjectType;
import graphql.schema.idl.SchemaDirectiveWiring;
import graphql.schema.idl.SchemaDirectiveWiringEnvironment;

/**
 * The authenticated directive for graphql queries.
 */
@DgsDirective(name = "authenticated")
public class AuthDirective implements SchemaDirectiveWiring {

  @Override
  public GraphQLFieldDefinition onField(
          SchemaDirectiveWiringEnvironment<GraphQLFieldDefinition> env) {

    if (!(env.getFieldsContainer() instanceof GraphQLObjectType container)) {
      return env.getElement();
    }

    GraphQLFieldDefinition field = env.getFieldDefinition();

    DataFetcher<?> original = env.getCodeRegistry()
            .getDataFetcher(container, field);

    if (original instanceof AuthenticatedFetcher) {
      return field;
    }

    GraphQLAppliedDirective dir = field.getAppliedDirective("authenticated");
    boolean requiresUser = getArgBool(dir, "requiresUser", true);
    boolean requiresOrg = getArgBool(dir, "requiresOrg", true);

    DataFetcher<?> wrapped = new AuthenticatedFetcher(original, requiresUser, requiresOrg);
    env.getCodeRegistry().dataFetcher(container, field, wrapped);

    return field;
  }

  /**
   * Gets the boolean value of the argument.
   *
   * @param dir the applied directive.
   * @param name the argument name.
   * @param def the default value.
   * @return the boolean value of the argument.
   */
  private boolean getArgBool(GraphQLAppliedDirective dir, String name, boolean def) {
    if (dir == null) {
      return def;
    }

    GraphQLAppliedDirectiveArgument arg = dir.getArgument(name);
    return arg != null && arg.getValue() != null ? (Boolean) arg.getValue() : def;
  }
}