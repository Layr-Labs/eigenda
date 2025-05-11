const fs = require("fs");
const path = require("path");

const goPath = path.resolve(__dirname, "wasm_exec.js");
require(goPath); // defines global `Go`

const go = new Go();

// Inject environment variables
go.env = {
  ...process.env, // include your shell's env
};

go.argv = ["relay.wasm", ...process.argv.slice(2)];

const wasmPath = path.resolve(__dirname, "relay.wasm");

(async () => {
  const wasmBuffer = fs.readFileSync(wasmPath);
  const { instance } = await WebAssembly.instantiate(wasmBuffer, go.importObject);
  try {
    await go.run(instance);
  } catch (err) {
    console.error("WASM runtime error:", err);
  }
})();
