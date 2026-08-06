import type * as React from "react";

import {
	DEFAULT_DESCRIPTION,
	OG_IMAGE_HEIGHT,
	OG_IMAGE_URL,
	OG_IMAGE_WIDTH,
	SITE_NAME,
	SITE_URL,
} from "#/constants";
import type { PageMeta } from "#/types";

// React's MetaHTMLAttributes has no `property`, which OpenGraph needs.
type MetaTag = React.JSX.IntrinsicElements["meta"] & { property?: string };
type LinkTag = React.JSX.IntrinsicElements["link"];

const absolute = (pathOrUrl: string) =>
	/^https?:\/\//.test(pathOrUrl)
 		? pathOrUrl
 		: `${SITE_URL}${pathOrUrl.startsWith("/") ? "" : "/"}${pathOrUrl}`;

/**
 * Builds the `head()` payload for a route: title, description, robots,
 * OpenGraph, Twitter card and canonical link.
 *
 *   head: () => seo(PAGE_META["/auth/login"])
 */
export function seo(page: PageMeta): { meta: MetaTag[]; links: LinkTag[] } {
	const description = page.description ?? DEFAULT_DESCRIPTION;
	const image = page.image ? absolute(page.image) : OG_IMAGE_URL;
	const url = page.path ? absolute(page.path) : undefined;

	const meta: MetaTag[] = [
		{ title: page.title },
		{ name: "description", content: description },
		{ property: "og:type", content: page.type ?? "website" },
		{ property: "og:site_name", content: SITE_NAME },
		{ property: "og:title", content: page.title },
		{ property: "og:description", content: description },
		{ property: "og:image", content: image },
		{ property: "og:image:width", content: String(OG_IMAGE_WIDTH) },
		{ property: "og:image:height", content: String(OG_IMAGE_HEIGHT) },
		{ name: "twitter:card", content: "summary_large_image" },
		{ name: "twitter:title", content: page.title },
		{ name: "twitter:description", content: description },
		{ name: "twitter:image", content: image },
	];

	if (page.robots) meta.push({ name: "robots", content: page.robots });
	if (url) meta.push({ property: "og:url", content: url });

	return { meta, links: url ? [{ rel: "canonical", href: url }] : [] };
}
