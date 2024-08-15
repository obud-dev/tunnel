import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "~": "/src",
    },
  },
  server: {
    proxy: {
      "/api": "http://0.0.0.0:8000",
    },
  },
});
