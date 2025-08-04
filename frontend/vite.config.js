import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss()],
  root: ".", // current dir is frontend root
  build: {
    outDir: "../public", // where Go will serve from
    emptyOutDir: true,
    assetsDir: "", // flat asset output (optional)
    rollupOptions: {
      output: {
        entryFileNames: "main.js",
        assetFileNames: "style.css",
      },
    },
  },
});
