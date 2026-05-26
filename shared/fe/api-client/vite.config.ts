import path from "node:path";
import {defineConfig} from "vite";

export default defineConfig({
  build: {
    minify: "oxc",
    reportCompressedSize: true,
    lib: {
      entry: path.resolve(__dirname, "src/index.ts"),
      name: "ApiClient",
      fileName: "index",
      formats: ["es"],
    },
    rollupOptions: {
      external: (id) => !id.startsWith(".") && !path.isAbsolute(id),
    },
  },
});
