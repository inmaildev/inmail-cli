package cmd

import (
	"fmt"
	"os"

	"github.com/inmaildev/inmail-cli/internal/client"
	"github.com/inmaildev/inmail-cli/internal/config"
	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfg                *config.Config
	apiClient          *client.Client
	flagAPIKey         string
	flagBaseURL        string
	flagOutput         string
	flagNonInteractive bool
)

var rootCmd = &cobra.Command{
	Use:   "inmail",
	Short: "InMail CLI — manage your email infrastructure from the terminal",
	Long: `InMail CLI provides programmatic access to the InMail API.

It supports both interactive (REPL) and non-interactive modes, making it
suitable for humans and AI agents alike.

Authentication:
  Set INMAIL_API_KEY environment variable, or pass --api-key, or run:
    inmail configure

Documentation:
  https://inmail.dev/docs

Examples:
  inmail messages list
  inmail users list --output json
  inmail send --to user@example.com --subject "Hello" --body "Hi there"
  inmail repl`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		output.SetFormat(flagOutput)

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		if flagAPIKey != "" {
			cfg.APIKey = flagAPIKey
		}
		if flagBaseURL != "" {
			cfg.BaseURL = flagBaseURL
		}

		skipAuthCmds := map[string]bool{
			"configure":  true,
			"version":    true,
			"help":       true,
			"completion": true,
		}
		if !skipAuthCmds[cmd.Name()] && cmd.Parent() != nil && !skipAuthCmds[cmd.Parent().Name()] {
			if cfg.APIKey == "" {
				return fmt.Errorf("no API key found. Set INMAIL_API_KEY or run: inmail configure")
			}
		}

		apiClient = client.New(cfg.BaseURL, cfg.APIKey)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "InMail API key (overrides INMAIL_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "API base URL (default: https://inmail.dev/v1)")
	rootCmd.PersistentFlags().StringVar(&flagOutput, "output", "human", "Output format: human or json")
	rootCmd.PersistentFlags().BoolVar(&flagNonInteractive, "non-interactive", false, "Disable prompts (for agent/CI use)")

	rootCmd.AddCommand(
		newConfigureCmd(),
		newVersionCmd(),
		newUsersCmd(),
		newMessagesCmd(),
		newAccountsCmd(),
		newConfigCmd(),
		newStatsCmd(),
		newSendCmd(),
		newSendLogsCmd(),
		newReplCmd(),
	)
}
