import { Spinner } from "../ui/Spinner";

export function RoutingPending() {
	return (
		<div className="bg-background flex min-h-screen items-center justify-center">
			<Spinner size="md" />
		</div>
	);
}
