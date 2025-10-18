import { graphql } from "../../helpers/graphql_request";
import {
  AircraftStatus,
  CreateAircraftDocument,
  CreateAircraftMutation,
  CreateAircraftMutationVariables,
  GetAircraftDocument,
  GetAircraftQuery,
  GetAircraftQueryVariables,
} from "../../gql/graphql";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check } from "k6";
import { AuthContext } from "src/types/auth_context";

export function runAircraftScenario(
  url: string,
  accessToken?: string,
  authContext?: AuthContext,
) {
  graphql<GetAircraftQuery, GetAircraftQueryVariables>(
    url,
    GetAircraftDocument,
    { id: uuidv4() },
    accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
  );

  const createAircraftRes = graphql<
    CreateAircraftMutation,
    CreateAircraftMutationVariables
  >(
    url,
    CreateAircraftDocument,
    {
      input: {
        registration: `N-${uuidv4().slice(0, 8)}`,
        manufacturer: "Boeing",
        model: "737",
        capacity: 50,
        status: AircraftStatus.Available,
        yearOfManufacture: 2020,
      },
    },
    {
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      ...(authContext
        ? {
            "x-org-id": authContext.organization.id,
            "x-user-sub": authContext.admin.id,
          }
        : {}),
    },
  );

  check(createAircraftRes, {
    "aircraft created": (r) => !!r?.createAircraft?.id,
  });

  const aircraftId = createAircraftRes?.createAircraft?.id;

  if (aircraftId) {
    graphql<GetAircraftQuery, GetAircraftQueryVariables>(
      url,
      GetAircraftDocument,
      { id: aircraftId },
      accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
    );
  }

  return aircraftId;
}
