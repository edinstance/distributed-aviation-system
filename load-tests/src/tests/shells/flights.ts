import { graphql } from "../../helpers/graphql_request";
import {
  CreateFlightDocument,
  CreateFlightMutation,
  CreateFlightMutationVariables,
  GetFlightByIdDocument,
  GetFlightByIdQuery,
  GetFlightByIdQueryVariables,
} from "../../gql/graphql";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check } from "k6";

function randomFlightNumber(): string {
  const airlines = ["BA", "UA", "LH", "AF"];
  const code = airlines[Math.floor(Math.random() * airlines.length)];
  return code + Math.floor(100 + Math.random() * 9000);
}

export function runFlightScenario(
  url: string,
  aircraftId: string,
  accessToken?: string,
) {
  graphql<GetFlightByIdQuery, GetFlightByIdQueryVariables>(
    url,
    GetFlightByIdDocument,
    { id: uuidv4() },
    accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
  );

  const createRes = graphql<
    CreateFlightMutation,
    CreateFlightMutationVariables
  >(
    url,
    CreateFlightDocument,
    {
      aircraftId,
      number: randomFlightNumber(),
      origin: "LHR",
      destination: "JFK",
      departureTime: new Date().toISOString(),
      arrivalTime: new Date(Date.now() + 3600 * 1000).toISOString(),
    },
    accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
  );

  check(createRes, {
    "flight created": (r) => !!r?.createFlight?.id,
  });

  const flightId = createRes?.createFlight?.id;

  if (flightId) {
    graphql<GetFlightByIdQuery, GetFlightByIdQueryVariables>(
      url,
      GetFlightByIdDocument,
      { id: flightId },
      accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
    );
  }

  return flightId;
}
