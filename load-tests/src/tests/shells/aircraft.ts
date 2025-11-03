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
import { buildAuthHeaders } from "../../helpers/auth_headers";

export function runAircraftScenario(
  url: string,
  accessToken?: string,
  authContext?: AuthContext,
) {
  // Determine if we should use direct headers (for router) vs JWT (for gateway)
  const useDirectHeaders = !url.includes("/gateway");

  graphql<GetAircraftQuery, GetAircraftQueryVariables>(
    url,
    GetAircraftDocument,
    { id: uuidv4() },
    buildAuthHeaders(accessToken, authContext, useDirectHeaders),
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
    buildAuthHeaders(accessToken, authContext, useDirectHeaders),
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
      buildAuthHeaders(accessToken, authContext, useDirectHeaders),
    );
  }

  return aircraftId;
}
