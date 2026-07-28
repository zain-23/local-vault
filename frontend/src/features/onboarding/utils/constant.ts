export const PLATFORMS = [
  { key: "macos", label: "macOS", cmd: "brew install localvault/tap/lv" },
  {
    key: "linux",
    label: "Linux",
    cmd: "curl -fsSL https://localvault.app/install.sh | sh",
  },
  {
    key: "source",
    label: "From source",
    cmd: "go install github.com/localvault/cli@latest",
  },
] as const;
