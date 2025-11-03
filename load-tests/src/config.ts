export const CONFIG = {
  flightServiceUrl: __ENV.FLIGHT_URL || "http://api.aviation.local/flights/graphql",
  aircraftServiceUrl: __ENV.AIRCRAFT_URL || "http://api.aviation.local/aircraft/graphql",
  routerUrl: __ENV.ROUTER_URL || "http://api.aviation.local/router",
  gatewayUrl: __ENV.GATEWAY_URL || "http://api.aviation.local/gateway",
  authServiceUrl: __ENV.AUTHENTICATION_URL || "http://api.aviation.local/auth",
};
