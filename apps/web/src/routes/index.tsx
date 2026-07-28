import { createFileRoute, redirect } from "@tanstack/react-router";

// Root lands on the members page for now (the first built protected screen).
export const Route = createFileRoute("/")({
	beforeLoad: () => {
		throw redirect({ to: "/members" });
	},
});
