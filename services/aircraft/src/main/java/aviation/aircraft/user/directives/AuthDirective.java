package aviation.aircraft.user.directives;

import aviation.aircraft.user.context.UserContext;
import com.netflix.graphql.dgs.DgsDirective;
import com.netflix.graphql.dgs.context.DgsContext;
import graphql.GraphQLException;
import graphql.schema.DataFetcher;
import graphql.schema.DataFetcherFactories;
import graphql.schema.GraphQLFieldDefinition;
import graphql.schema.GraphQLFieldsContainer;
import graphql.schema.GraphQLObjectType;
import graphql.schema.idl.SchemaDirectiveWiring;
import graphql.schema.idl.SchemaDirectiveWiringEnvironment;

/**
 * A directive for authentication.
 */
@DgsDirective(name = "authentication")
public class AuthDirective implements SchemaDirectiveWiring {

  /**
   * The directive implementation.
   *
   * @param env the environment for the directive.
   * @return the field definition with the authentication.
   */
  @Override
  public GraphQLFieldDefinition onField(SchemaDirectiveWiringEnvironment
                                                  <GraphQLFieldDefinition> env) {
    GraphQLFieldsContainer fieldsContainer = env.getFieldsContainer();
    GraphQLFieldDefinition fieldDefinition = env.getFieldDefinition();

    DataFetcher<?> originalDataFetcher = env.getCodeRegistry()
            .getDataFetcher((GraphQLObjectType) fieldsContainer, fieldDefinition);

    DataFetcher<?> authFetcher = DataFetcherFactories.wrapDataFetcher(
            originalDataFetcher,
            (dataFetchingEnvironment, value) -> {
              UserContext user = DgsContext.getCustomContext(dataFetchingEnvironment);

              if (user == null || user.getUserId() == null || user.getOrgId() == null) {
                throw new GraphQLException("Unauthorized: userSub or orgId missing");
              }

              try {
                return originalDataFetcher.get(dataFetchingEnvironment);
              } catch (Exception e) {
                throw new RuntimeException(e);
              }

            }
    );

    env.getCodeRegistry().dataFetcher((GraphQLObjectType) fieldsContainer,
            fieldDefinition, authFetcher);

    return fieldDefinition;
  }
}