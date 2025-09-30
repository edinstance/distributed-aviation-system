import gql from "graphql-tag";

export const GET_FLIGHT_BY_ID_QUERY = gql`
  query GetFlightById($id: ID!) {
    getFlightById(id: $id) {
      id
      number
      origin
      destination
      status
      aircraft {
        id
      }
    }
  }
`;

export const GET_FLIGHT_BY_ID_WITH_AIRCRAFT_QUERY = gql`
  query GetFlightByIdWithAircraft($id: ID!) {
    getFlightById(id: $id) {
      id
      number
      origin
      destination
      status
      aircraft {
        id
        registration
        manufacturer
        model
        status
      }
    }
  }
`;
