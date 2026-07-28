<div align="center">

<br/>

# 🔐 LocalVault

**Encrypted secret sync for dev teams.**
**No cloud. No leaks. Pure P2P.**

> Replace `.env` files with an encrypted vault that syncs directly
> between developer machines. Your secrets never touch a server in readable form.

<br/>

</div>

---

## Install

Prebuilt binaries are published to [GitHub Releases](https://github.com/zain-23/local-vault/releases) for Linux, macOS, and Windows.

**Linux / macOS** — download the right archive, extract, and install:

```bash
# pick your platform: linux_amd64 · linux_arm64 · darwin_amd64 · darwin_arm64
curl -fsSL https://github.com/zain-23/local-vault/releases/latest/download/lv_linux_amd64.tar.gz \
  | tar -xz lv
sudo mv lv /usr/local/bin/lv
```

**Windows** — download `lv_windows_amd64.zip` from the
[latest release](https://github.com/zain-23/local-vault/releases/latest),
unzip, and add `lv.exe` to your `PATH`.

**With Go**

```bash
go install github.com/zain-23/local-vault@latest
```

**From source**

```bash
git clone https://github.com/zain-23/local-vault
cd local-vault
go build -o lv ./apps/cli
# or: task cli:build
sudo mv lv /usr/local/bin/lv
```

> The release binaries default to the hosted server. To point `lv` at a
> different server, set `SERVER_URL`, e.g. `export SERVER_URL=https://your-server`.

---

## Quick Start

```bash
# Initialize vault in your project
cd my-project
lv init

# Unlock once per day — no prompts after this
lv unlock

# Add secrets
lv add DATABASE_URL=postgres://localhost/mydb
lv add API_KEY=sk-live-xxx
lv add STRIPE_KEY=sk_live_xxx --env production
lv add STRIPE_KEY=sk_test_xxx --env development

# Run your project with secrets injected
lv inject -- npm run dev

# Import existing .env file
lv import .env.local
```

---

## Session Management

LocalVault asks for your passphrase **once per day**. After that every command works without any prompt — like SSH agent.

```bash
# Unlock once
lv unlock
# ✅ Unlocked for 12 hours

# All commands work silently
lv list
lv add KEY=value
lv push
lv sync

# Lock when done
lv lock

# Lock all projects at once
lv lock --all
```

> Each project has its own independent session. Unlocking one project never affects another.

---

## Team Sync

```bash
# ── Dev A ──────────────────────────────────────────

lv invite
# ⏳ Waiting for teammate to join...
# 🎉 Teammate joined!
# Now run: lv push

lv push
# ✅ Sent to ahmed-laptop


# ── Dev B ──────────────────────────────────────────

lv join LV-A3F9-X2K1
# ✅ Connected to peer

lv sync
# ✅ Received 3 secret(s)

lv inject -- npm run dev  # ✅ all secrets loaded
```

**Works from anywhere.** Different offices, different cities, different networks — the signaling server bridges peers securely. If Dev B is offline when Dev A pushes, messages are stored for 48 hours and delivered when they come online.

---

## Secret Rotation

```bash
# Rotate single key
lv rotate API_KEY
# New value: ••••••••
# ✅ API_KEY rotated
# 📤 All peers notified

# Rotate multiple keys
lv rotate API_KEY STRIPE_KEY DATABASE_URL

# Open all secrets in editor — change only what you want
lv rotate --all

# Rotate specific environment
lv rotate --all --env production
```

---

## Peer Management

```bash
# See who has access
lv peers
# 👥 Trusted Peers
#   Peer 1 — ahmed-laptop   (b427d8d8...)
#   Peer 2 — sara-macbook   (c891f2a1...)

# Remove a peer
lv revoke b427d8d8-8b7f-4605-9b5a-f8be956b168d
# ✅ Revoked access for ahmed-laptop

# Rotate after revoking so cached copy is useless
lv rotate --all
```

---

## Security

```
Passphrase → Argon2id → AES-256-GCM → vault.json.enc
```

```
X25519 Key Exchange → Shared Secret → AES-256-GCM → Encrypted blob
```

| Layer             | Algorithm   |
| ----------------- | ----------- |
| Vault encryption  | AES-256-GCM |
| Key derivation    | Argon2id    |
| Peer key exchange | X25519 ECDH |
| Device identity   | Ed25519     |
| Session cache     | OS Keychain |

**What the signaling server sees:**

- Encrypted blobs — cannot read
- Random device IDs — not linked to identity
- IP addresses — only for peer discovery

**What never leaves your machine unencrypted:**

- Your secrets
- Your passphrase
- Your private keys

---

## All Commands

| Command                             | Description                   |
| ----------------------------------- | ----------------------------- |
| `lv init`                           | Initialize encrypted vault    |
| `lv unlock`                         | Unlock for 12 hours           |
| `lv lock`                           | Lock current project          |
| `lv lock --all`                     | Lock all projects             |
| `lv add KEY=VALUE`                  | Add or update a secret        |
| `lv add KEY=VALUE --env production` | Add to specific environment   |
| `lv get KEY`                        | Get value of a secret         |
| `lv list`                           | List all secrets              |
| `lv list --env staging`             | Filter by environment         |
| `lv remove KEY`                     | Delete a secret               |
| `lv import FILE`                    | Import from `.env` file       |
| `lv inject -- CMD`                  | Run command with secrets      |
| `lv inject --env production -- CMD` | Run with environment secrets  |
| `lv invite`                         | Generate teammate invite code |
| `lv join CODE`                      | Join teammate vault           |
| `lv push`                           | Push secrets to all peers     |
| `lv sync`                           | Pull secrets from peers       |
| `lv peers`                          | List trusted peers            |
| `lv rotate KEY`                     | Rotate a secret               |
| `lv rotate --all`                   | Rotate all in editor          |
| `lv revoke DEVICE_ID`               | Remove peer access            |
| `lv status`                         | Show vault health             |
| `lv log`                            | Show audit trail              |
| `lv whoami`                         | Show device identity          |

---

## Next.js

```bash
lv import .env.local --env development
lv import .env.production --env production
```

```json
{
  "scripts": {
    "dev": "lv inject --env development -- next dev",
    "build": "lv inject --env production  -- next build",
    "start": "lv inject --env production  -- next start"
  }
}
```

Your app code stays exactly the same — `process.env.KEY` works as always.

---

## Monorepo development

This repo is a pnpm + Turborepo monorepo with Go apps orchestrated via [Task](https://taskfile.dev):

```text
apps/cli      Go CLI (`lv`)
apps/server   Go API + worker
apps/web      React web app
packages/     Shared TypeScript packages
```

```bash
pnpm install          # JS workspaces
pnpm dev              # web (Turbo)
pnpm build            # all JS packages
pnpm lint / test      # JS lint & tests

task cli:build        # build lv
task server:run       # run API
task test:go          # go test ./...
task docker:server    # Docker image for API
```

Requires [pnpm](https://pnpm.io), [Go](https://go.dev), and optionally [Task](https://taskfile.dev).

---

## What Gets Committed to Git

```
.lv/identity.pub     ✅  safe to commit
.lv/vault.json.enc   ❌  gitignored automatically
.lv/identity.key     ❌  gitignored automatically
.lv/identity.json    ❌  gitignored automatically
```

---

<div align="center">
<br/>

_Stop sharing secrets over Slack._

[Report Bug](https://github.com/zain-23/local-vault/issues) · [Request Feature](https://github.com/zain-23/local-vault/issues)

</div>
