package cmd

import (
	"fmt"
	"os"
	"os/exec" // runs external commands (like child_process in Node.js)

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

var injectCmd = &cobra.Command{
	Use:   "inject [-- command]",
	Short: "Inject secrets into shell or run a command with secrets",
	Example: `  lv inject                    # print export statements
  lv inject -- npm run dev     # run command with secrets injected
  lv inject --env staging -- npm start`,
	RunE: func(cmd *cobra.Command, args []string) error {
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		dir, _ := os.Getwd()
		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Get secrets as key→value map
		secrets := v.InjectMap(envFlag)

		// If no command provided, print export statements
		// User can eval this: eval $(lv inject)
		if len(args) == 0 {
			for key, value := range secrets {
				fmt.Printf("export %s=%q\n", key, value)
			}
			return nil
		}

		// Run the provided command with secrets injected
		// args[0] is the command, args[1:] are its arguments
		// Example: args = ["npm", "run", "dev"]
		command := exec.Command(args[0], args[1:]...)

		// Start with current environment
		// Then add our secrets on top
		// os.Environ() = all current env vars (like process.env in Node.js)
		command.Env = os.Environ()
		for key, value := range secrets {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
		}

		// Connect command's stdin/stdout/stderr to our terminal
		// So the spawned process feels native (colors, interactive input, etc.)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		fmt.Printf("🚀 Injected %d secrets, starting: %s\n\n", len(secrets), args[0])

		// Run the command and wait for it to finish
		// Like: child_process.spawnSync() in Node.js
		if err := command.Run(); err != nil {
			// Exit with same code as the child process
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return err
		}

		return nil
	},
}

func init() {
	injectCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment to inject")
	rootCmd.AddCommand(injectCmd)
}
