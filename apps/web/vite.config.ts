import tailwindcss from "@tailwindcss/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { nitro } from "nitro/vite";
import { defineConfig, loadEnv } from "vite";

const config = defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), "VITE_");
	if (mode === "production" && !env.VITE_API_URL) {
		throw new Error("Missing VITE_API_URL");
	}

	return {
		resolve: { tsconfigPaths: true },
		plugins: [devtools(), tailwindcss(), nitro(), tanstackStart(), viteReact()],
	};
});

export default config;
