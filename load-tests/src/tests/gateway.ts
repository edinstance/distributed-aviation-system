import { sleep } from "k6";
import { Options } from "k6/options";

import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";
import { runAuthenticationScenario } from "./shells/auth";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export default function () {
  const accessToken = runAuthenticationScenario(CONFIG.authServiceUrl);

  const validAccessToken =
    typeof accessToken === "string" ? accessToken : undefined;

  const aircraftId = runAircraftScenario(CONFIG.gatewayUrl, validAccessToken);
  if (aircraftId) {
    runFlightScenario(CONFIG.gatewayUrl, aircraftId, validAccessToken);
  }
  sleep(1);
}
