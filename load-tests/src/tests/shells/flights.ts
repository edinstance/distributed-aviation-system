import { graphql } from "../../helpers/graphql_request";
import {
  CreateFlightDocument,
  CreateFlightMutation,
  CreateFlightMutationVariables,
  GetFlightByIdDocument,
  GetFlightByIdQuery,
  GetFlightByIdQueryVariables,
} from "../../gql/graphql";
import { CONFIG } from "../../config";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import { check } from "k6";

function randomFlightNumber(): string {
  const airlines = ["BA", "UA", "LH", "AF"];
  const code = airlines[Math.floor(Math.random() * airlines.length)];
  return code + Math.floor(100 + Math.random() * 9000);
}

export function runFlightScenario(aircraftId: string) {
  graphql<GetFlightByIdQuery, GetFlightByIdQueryVariables>(
    CONFIG.flightServiceUrl,
    GetFlightByIdDocument,
    { id: uuidv4() },
  );

  const createRes = graphql<
    CreateFlightMutation,
    CreateFlightMutationVariables
  >(CONFIG.flightServiceUrl, CreateFlightDocument, {
    aircraftId,
    number: randomFlightNumber(),
    origin: "LHR",
    destination: "JFK",
    departureTime: new Date().toISOString(),
    arrivalTime: new Date(Date.now() + 3600 * 1000).toISOString(),
  });

  check(createRes, {
    "flight created": (r) => !!r?.createFlight?.id,
  });

  const flightId = createRes?.createFlight?.id;

  if (flightId) {
    graphql<GetFlightByIdQuery, GetFlightByIdQueryVariables>(
      CONFIG.flightServiceUrl,
      GetFlightByIdDocument,
      { id: flightId },
    );
  }

  return flightId;
}
