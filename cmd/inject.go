package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject [-- command]",
	Short: "Inject secrets into shell or run a command with secrets",
	Example: `  lv inject                    # print export statements
  lv inject -- npm run dev     # run command with secrets injected
  lv inject --env staging -- npm start`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := os.Getwd()

		// Use session key — no passphrase needed
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		// Get secrets as key→value map
		secrets := v.InjectMap(envFlag)

		if len(secrets) == 0 {
			fmt.Println("No secrets found for this environment.")
			fmt.Println("Add secrets with: lv add KEY=value")
			return nil
		}

		// No command provided — print export statements
		// User can eval: eval $(lv inject)
		if len(args) == 0 {
			for key, value := range secrets {
				fmt.Printf("export %s=%q\n", key, value)
			}
			return nil
		}

		// Run command with secrets injected as env vars
		command := exec.Command(args[0], args[1:]...)

		// Start with current shell env
		// Then add vault secrets on top
		command.Env = os.Environ()
		for key, value := range secrets {
			command.Env = append(command.Env,
				fmt.Sprintf("%s=%s", key, value))
		}

		// Connect to terminal so spawned process feels native
		// Colors, interactive input, live output all work correctly
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		fmt.Printf("🚀 Injected %d secrets, starting: %s\n\n",
			len(secrets), args[0])

		// Run and wait for process to finish
		if err := command.Run(); err != nil {
			// Exit with same code as child process
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
