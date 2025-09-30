import { sleep } from "k6";
import { Options } from "k6/options";

import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export default function () {
  const aircraftId = runAircraftScenario(CONFIG.routerUrl);
  if (aircraftId) {
    runFlightScenario(CONFIG.routerUrl, aircraftId);
  }
  sleep(1);
}
