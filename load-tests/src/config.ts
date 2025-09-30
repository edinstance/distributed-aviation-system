export const CONFIG = {
  flightServiceUrl: __ENV.FLIGHT_URL || "http://localhost:8081/graphql",
  aircraftServiceUrl: __ENV.AIRCRAFT_URL || "http://localhost:8080/graphql",
  routerUrl: __ENV.ROUTER_URL || "http://localhost:4000",
};
