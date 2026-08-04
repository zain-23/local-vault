import { createFileRoute, redirect } from "@tanstack/react-router";

// Root lands on the workspace dashboard.
export const Route = createFileRoute("/")({
	beforeLoad: () => {
		throw redirect({ to: "/dashboard" });
	},
});
