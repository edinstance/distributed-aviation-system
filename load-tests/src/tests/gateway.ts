import { sleep } from "k6";
import { Options } from "k6/options";

import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";
import { runAuthenticationScenario } from "./shells/auth";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export function setup() {
  const accessToken = runAuthenticationScenario(CONFIG.authServiceUrl);

  const validAccessToken =
    typeof accessToken === "string" ? accessToken : undefined;

  const aircraftId = runAircraftScenario(CONFIG.gatewayUrl, validAccessToken);

  return { aircraftId, validAccessToken };
}

export default function (data: { aircraftId: string; validAccessToken?: string }) {
  if (data.aircraftId && data.validAccessToken) {
    runFlightScenario(CONFIG.gatewayUrl, data.aircraftId, data.validAccessToken);
  }
  sleep(1);
}
