import { resolve } from "path";
import { defineConfig } from "vite";

export default defineConfig({
  server: {
    port: 3000,
  },
  build: {
    target: "esnext",
    rollupOptions: {
      input: {
        main: resolve(__dirname, "index.html"),
        "mandelbrot-webgl": resolve(__dirname, "mandelbrot-webgl/index.html"),
        "kaboom-webgl": resolve(__dirname, "kaboom-webgl/index.html"),
      },
    },
  },
});
