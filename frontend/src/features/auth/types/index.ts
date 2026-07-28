import type * as React from "react";

// One line of the brand terminal session: a typed command (`cmd`) or its
// output (`out`). `id` gives each line a stable React key.
export type TerminalEntry = { id: string } & (
	| { cmd: string }
	| { out: React.ReactNode }
);
