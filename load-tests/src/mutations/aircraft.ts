import gql from "graphql-tag";

export const CREATE_AIRCRAFT_MUTATION = gql`
  mutation CreateAircraft($input: CreateAircraftInput!) {
    createAircraft(input: $input) {
      id
      registration
      manufacturer
      model
      status
    }
  }
`;
