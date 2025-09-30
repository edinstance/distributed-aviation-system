import gql from "graphql-tag";

export const CREATE_FLIGHT_MUTATION = gql`
  mutation CreateFlight(
    $number: String!
    $origin: String!
    $destination: String!
    $departureTime: Time!
    $arrivalTime: Time!
    $aircraftId: ID!
  ) {
    createFlight(
      number: $number
      origin: $origin
      destination: $destination
      departureTime: $departureTime
      arrivalTime: $arrivalTime
      aircraftId: $aircraftId
    ) {
      id
      number
      origin
      destination
      departureTime
      arrivalTime
      aircraft {
        id
      }
    }
  }
`;
