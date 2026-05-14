# web/

Frontend for LocalVault. Lives outside the Go module — `go build ./...` will
not descend into this directory.

When scaffolding, pick one of:

- **Next.js** — if you need SSR / API routes
- **Vite + React** — pure SPA talking to `lv-server`
- **SvelteKit** — lighter alternative to Next

The frontend talks to `lv-server` over HTTP. Route definitions live in
`server/main.go` (and will be split into `server/handlers/` as the API grows).
