import { Avatar, AvatarFallback, AvatarImage } from "#/components/ui";
import type { AuditEvent } from "#/features/audit/api";
import {
	ACTION_TONE_CLASS,
	actionLabel,
	actionTone,
	formatEventClock,
	hueFromString,
	resolveTargetLabel,
} from "#/features/audit/utils";
import { cn } from "#/lib/utils.ts";
import { initialsFromName } from "#/utils";


type AuditEventRowProps = {
	event: AuditEvent;
};

export function AuditEventRow({ event }: AuditEventRowProps) {
	const actor = event.actor_name?.trim() || "Someone";
	const verb = actionLabel(event.action);
	const tone = actionTone(event.action);
	const target = resolveTargetLabel(event);
	const meta = [event.device_id, event.ip].filter(Boolean).join(" · ");
	const hue = hueFromString(event.actor_id || actor);
	const initials = initialsFromName(actor);

	return (
		<div className="grid grid-cols-[64px_36px_minmax(0,1fr)] items-start gap-3 border-b border-border py-2.5 last:border-b-0 sm:grid-cols-[76px_36px_minmax(0,1fr)] sm:gap-3.5 sm:px-1">
			<time
				dateTime={event.created_at}
				className="pt-1 font-mono text-xs text-muted-foreground tabular-nums"
			>
				{formatEventClock(event.created_at)}
			</time>

			<Avatar aria-hidden className="rounded-lg">
				{event.actor_avatar_url ? (
					<AvatarImage src={event.actor_avatar_url} alt={actor} />
				) : null}
				<AvatarFallback
					className="rounded-lg text-sm font-semibold text-white"
					style={{ background: `oklch(0.52 0.12 ${hue})` }}
				>
					{initials}
				</AvatarFallback>
			</Avatar>

			<div className="min-w-0">
				<p className="text-sm leading-snug">
					<span className="font-medium text-foreground">{actor}</span>
					<span className={cn("mx-1.5 font-medium", ACTION_TONE_CLASS[tone])}>
						{verb}
					</span>
					{target ? (
						<span className="font-mono text-[12.5px] text-foreground/90">
							{target}
						</span>
					) : null}
				</p>
				{meta ? (
					<p className="mt-0.5 font-mono text-[11.5px] text-muted-foreground">
						{meta}
					</p>
				) : (
					<p className="mt-0.5 font-mono text-[11.5px] text-muted-foreground">
						{event.action}
					</p>
				)}
			</div>
		</div>
	);
}
