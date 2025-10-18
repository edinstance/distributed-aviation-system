import { sleep } from "k6";
import { Options } from "k6/options";

import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";
import { createOrganization } from "src/helpers/organization";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export default function () {
  const authContext = createOrganization(CONFIG.authServiceUrl);

  const aircraftId = runAircraftScenario(
    CONFIG.routerUrl,
    undefined,
    authContext,
  );
  if (aircraftId) {
    runFlightScenario(CONFIG.routerUrl, aircraftId, undefined, authContext);
  }
  sleep(1);
}
