import { LVLogo } from "#/components/shared/LVLogo.tsx";
import {
	FOOTER_BLURB,
	FOOTER_COLUMNS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { Container } from "./Container.tsx";

function LandingFooter() {
	return (
		<footer className="border-t border-border pt-12 pb-10">
			<Container>
				<div className="grid grid-cols-2 gap-8 md:grid-cols-[1.6fr_1fr_1fr]">
					<div>
						<a href={`#${SECTION_IDS.top}`} aria-label="LocalVault home">
							<LVLogo size={22} className="gap-2.5" />
						</a>
						<p className="mt-3 max-w-[34ch] text-[13.5px] leading-relaxed text-muted-foreground">
							{FOOTER_BLURB}
						</p>
					</div>

					{FOOTER_COLUMNS.map((column) => (
						<div key={column.id}>
							<h2 className="mb-3 text-[13px] font-medium">{column.title}</h2>
							{column.links.map((link) => (
								<a
									key={link.label}
									href={link.href}
									target={link.external ? "_blank" : undefined}
									rel={link.external ? "noreferrer noopener" : undefined}
									className="block py-1 text-xs text-muted-foreground transition-colors hover:text-foreground"
								>
									{link.label}
								</a>
							))}
						</div>
					))}
				</div>

				<div className="mt-10 flex flex-wrap justify-between gap-4 border-t border-border pt-5 text-[12.5px] text-muted-foreground">
					<span>© 2026 LocalVault — MIT licensed.</span>
					<span>No telemetry. No analytics. No account required.</span>
				</div>
			</Container>
		</footer>
	);
}

export { LandingFooter };
