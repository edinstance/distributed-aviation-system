import { sleep } from "k6";
import { Options } from "k6/options";

import { runFlightScenario } from "./shells/flights";
import { runAircraftScenario } from "./shells/aircraft";
import { CONFIG } from "../config";
import { createOrganization } from "src/helpers/organization";
import { AuthContext } from "src/types/auth_context";

export let options: Options = {
  stages: [{ duration: "10s", target: 100 },{ duration: "10s", target: 1000 } ],
};

export function setup() {
  const authContext = createOrganization(CONFIG.authServiceUrl);

  const aircraftId = runAircraftScenario(
    CONFIG.routerUrl,
    undefined,
    authContext,
  );

  return { aircraftId, authContext };
}

export default function (data: { aircraftId: string ,authContext: AuthContext }) {
  if (data.aircraftId && data.authContext) {
    runFlightScenario(
      CONFIG.routerUrl,
      data.aircraftId,
      undefined,
      data.authContext,
    );
  }
  sleep(1);
}
