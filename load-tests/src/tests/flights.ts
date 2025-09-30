import { sleep } from "k6";
import { Options } from "k6/options";
import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export function setup() {
  const aircraftId = runAircraftScenario(CONFIG.aircraftServiceUrl);
  return { aircraftId };
}

export default function (data: { aircraftId: string }) {
  if (data.aircraftId) {
    runFlightScenario(CONFIG.flightServiceUrl, data.aircraftId);
  }
  sleep(1);
}
