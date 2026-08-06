import { Button } from "#/components/ui";
import { cn } from "#/lib/utils.ts";

export type PaginationProps = {
	page: number;
	totalPages: number;
	onPageChange: (page: number) => void;
	className?: string;
	/** Hide when there's only one page (default true). */
	hideWhenSinglePage?: boolean;
};

// Prev/next pager for server-paginated lists. Domain-agnostic — pass page
// numbers, not API meta shapes, so any feature can reuse it.
export function Pagination({
	page,
	totalPages,
	onPageChange,
	className,
	hideWhenSinglePage = true,
}: PaginationProps) {
	if (hideWhenSinglePage && totalPages <= 1) return null;

	return (
		<div
			className={cn(
				"flex items-center justify-between text-sm text-muted-foreground",
				className,
			)}
		>
			<span>
				Page {page} of {totalPages}
			</span>
			<div className="flex gap-2">
				<Button
					variant="outline"
					size="sm"
					disabled={page <= 1}
					onClick={() => onPageChange(page - 1)}
				>
					Previous
				</Button>
				<Button
					variant="outline"
					size="sm"
					disabled={page >= totalPages}
					onClick={() => onPageChange(page + 1)}
				>
					Next
				</Button>
			</div>
		</div>
	);
}
