import { sleep } from "k6";
import { Options } from "k6/options";
import { CONFIG } from "../config";
import { runAuthenticaionScenario } from "./shells/auth";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};
export default function () {
  runAuthenticaionScenario(CONFIG.authServiceUrl);
  sleep(1);
}
