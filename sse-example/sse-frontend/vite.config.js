import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/data": "http://localhost:8080",
      "/events": "http://localhost:8080",
    },
  },
});