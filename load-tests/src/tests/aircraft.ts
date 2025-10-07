import { sleep } from "k6";
import { runAircraftScenario } from "./shells/aircraft";
import { Options } from "k6/options";
import { CONFIG } from "../config";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};
export default function () {
  runAircraftScenario(CONFIG.aircraftServiceUrl);
  sleep(1);
}
