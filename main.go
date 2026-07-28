package main

// Thin module-root entry so `go install github.com/zain-23/local-vault@latest`
// and `go build -o lv .` keep working after the CLI moved to apps/cli.

import "github.com/zain-23/local-vault/apps/cli/cmd"

func main() {
	cmd.Execute()
}
