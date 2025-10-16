import { build } from "esbuild";

await build({
  entryPoints: [
    "src/tests/flights.ts",
    "src/tests/aircraft.ts",
    "src/tests/router.ts",
    "src/tests/auth.ts",
  ],
  bundle: true,
  format: "esm",
  platform: "neutral",
  outdir: "dist",
  external: ["k6", "k6/*"],
});
