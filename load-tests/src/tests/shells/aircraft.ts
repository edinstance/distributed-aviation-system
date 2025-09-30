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
import { CONFIG } from "../../config";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check } from "k6";

export function runAircraftScenario() {
  graphql<GetAircraftQuery, GetAircraftQueryVariables>(
    CONFIG.aircraftServiceUrl,
    GetAircraftDocument,
    { id: uuidv4() },
  );

  const createAircraftRes = graphql<
    CreateAircraftMutation,
    CreateAircraftMutationVariables
  >(CONFIG.aircraftServiceUrl, CreateAircraftDocument, {
    input: {
      registration: `N${Math.floor(Math.random() * 90000) + 10000}`,
      manufacturer: "Boeing",
      model: "737",
      capacity: 50,
      status: AircraftStatus.Available,
      yearOfManufacture: 2020,
    },
  });

  check(createAircraftRes, {
    "aircraft created": (r) => !!r?.createAircraft?.id,
  });

  const aircraftId = createAircraftRes?.createAircraft?.id;

  if (aircraftId) {
    graphql<GetAircraftQuery, GetAircraftQueryVariables>(
      CONFIG.aircraftServiceUrl,
      GetAircraftDocument,
      { id: aircraftId },
    );
  }

  return aircraftId;
}
