import { sleep } from "k6";
import { runAircraftScenario } from "./shells/aircraft";
import { runFlightScenario } from "./shells/flights";
import { Options } from "k6/options";

export let options: Options = {
  stages: [{ duration: "10s", target: 1 }],
};

export default function () {
  const aircraftId = runAircraftScenario();
  if (aircraftId) {
    runFlightScenario(aircraftId);
  }
  sleep(1);
}
