import { sleep } from "k6";
import { runAircraftScenario } from "./shells/aircraft";
import { Options } from "k6/options";
import { CONFIG } from "../config";
import { createOrganization } from "src/helpers/organization";
import { AuthContext } from "src/types/auth_context";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export function setup() {
  const authContext = createOrganization(CONFIG.authServiceUrl);

  return { authContext };
}

export default function (data: { authContext: AuthContext }) {
  if (data.authContext) {
    runAircraftScenario(CONFIG.aircraftServiceUrl, undefined, data.authContext);
  }
  sleep(1);
}
