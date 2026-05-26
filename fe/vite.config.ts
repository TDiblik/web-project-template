import path from "node:path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import {defineConfig} from "vite";

export default defineConfig(() => ({
  plugins: [react(), tailwindcss()],
  build: {
    target: "baseline-widely-available",
    modulePreload: {
      polyfill: true,
    },
    outDir: "dist",
    minify: "oxc",
    cssMinify: "lightningcss",
    sourcemap: false,
    ssr: false,
    reportCompressedSize: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("shared/fe")) return "shared-fe";

          if (id.includes("node_modules")) {
            if (id.includes("motion")) return "motion-vendor";
            if (id.includes("@tanstack")) return "tanstack-vendor";
            if (id.includes("i18next")) return "i18n-vendor";
            return "vendor";
          }
        },
      },
    },
  },
  resolve: {
    alias: {
      "@shared/api-client": path.resolve(__dirname, "../shared/fe/api-client/dist/index"),
    },
  },
  server: {
    port: 35230,
  },
}));
