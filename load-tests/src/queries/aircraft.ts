import gql from "graphql-tag";

export const GET_AIRCRAFT_BY_ID = gql`
  query GetAircraft($id: ID!) {
    getAircraftById(input: $id) {
      id
      registration
      manufacturer
      model
      status
    }
  }
`;
