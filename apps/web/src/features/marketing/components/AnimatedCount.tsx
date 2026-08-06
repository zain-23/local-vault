import {
	animate,
	motion,
	useMotionValue,
	useReducedMotion,
	useTransform,
} from "motion/react";
import { useEffect } from "react";

const COUNT_DURATION = 0.9;
/** easeOutCubic — fast off the mark, long settle. */
const COUNT_EASE = [0.33, 1, 0.68, 1] as const;

/**
 * Ticks from zero to `value` on mount. The number lives in a motion value, so
 * the count runs off the animation frame loop without re-rendering React once.
 */
function AnimatedCount({
	value,
	className,
}: {
	value: number;
	className?: string;
}) {
	const prefersReducedMotion = useReducedMotion();
	const count = useMotionValue(0);
	const display = useTransform(count, (latest) =>
		Math.round(latest).toLocaleString(),
	);

	useEffect(() => {
		if (prefersReducedMotion) {
			count.set(value);
			return;
		}
		const controls = animate(count, value, {
			duration: COUNT_DURATION,
			ease: [...COUNT_EASE],
		});
		return () => controls.stop();
	}, [count, prefersReducedMotion, value]);

	return <motion.span className={className}>{display}</motion.span>;
}

export { AnimatedCount };
