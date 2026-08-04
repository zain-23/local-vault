import type { BreadcrumbStaticData } from "#/components/layout/types";

export function readBreadcrumb(staticData: unknown): string | undefined {
	const label = (staticData as BreadcrumbStaticData | undefined)?.breadcrumb;
	return typeof label === "string" ? label : undefined;
}
