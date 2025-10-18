import { sleep } from "k6";
import { runAircraftScenario } from "./shells/aircraft";
import { Options } from "k6/options";
import { CONFIG } from "../config";
import { createOrganization } from "src/helpers/organization";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};
export default function () {
  const organization = createOrganization(CONFIG.authServiceUrl);

  runAircraftScenario(CONFIG.aircraftServiceUrl, undefined, organization);
  sleep(1);
}
