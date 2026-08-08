import { useQuery } from "@tanstack/react-query";
import { GITHUB_REPO } from "#/constants";

async function fetchStarCount(): Promise<number> {
	const res = await fetch(`https://api.github.com/repos/${GITHUB_REPO}`);
	if (!res.ok) throw new Error("Failed to fetch GitHub repo stats");
	const data: { stargazers_count: number } = await res.json();
	return data.stargazers_count;
}

/** Star count for the nav badge. GitHub's public API needs no auth, so this
 * hits it directly rather than proxying through our own backend. */
export function useGithubStars() {
	return useQuery({
		queryKey: ["github-stars", GITHUB_REPO],
		queryFn: fetchStarCount,
		staleTime: 60 * 60 * 1000,
		retry: false,
	});
}
