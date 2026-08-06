import { Link } from "@tanstack/react-router";
import { Fragment } from "react";

import { useBreadcrumbs } from "#/components/layout/hooks";
import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbLink,
	BreadcrumbList,
	BreadcrumbPage,
	BreadcrumbSeparator,
} from "#/components/ui";

export function AppBreadcrumb() {
	const crumbs = useBreadcrumbs();
	if (crumbs.length === 0) return null;

	return (
		<Breadcrumb>
			<BreadcrumbList className="flex-nowrap">
				{crumbs.map((crumb, index) => {
					const isLast = index === crumbs.length - 1;
					return (
						<Fragment key={crumb.to}>
							{index > 0 && <BreadcrumbSeparator />}
							<BreadcrumbItem className="min-w-0">
								{isLast ? (
									<BreadcrumbPage className="truncate">
										{crumb.label}
									</BreadcrumbPage>
								) : (
									<BreadcrumbLink asChild className="truncate">
										<Link to={crumb.to}>{crumb.label}</Link>
									</BreadcrumbLink>
								)}
							</BreadcrumbItem>
						</Fragment>
					);
				})}
			</BreadcrumbList>
		</Breadcrumb>
	);
}
