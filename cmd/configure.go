package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/inmaildev/inmail-cli/internal/config"
	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newConfigureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Save API key and base URL to ~/.inmail/config.json",
		Long: `Interactively configure the InMail CLI.

Your API key is stored in ~/.inmail/config.json with 600 permissions.

Get your API key at: https://inmail.dev/admin/api-keys`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagNonInteractive {
				if flagAPIKey == "" {
					output.PrintError("--api-key is required in --non-interactive mode")
					return fmt.Errorf("missing --api-key")
				}
				baseURL := "https://inmail.dev/v1"
				if flagBaseURL != "" {
					baseURL = flagBaseURL
				}
				c := &config.Config{APIKey: flagAPIKey, BaseURL: baseURL}
				if err := config.Save(c); err != nil {
					output.PrintError(err.Error())
					return err
				}
				output.PrintSuccess("Configuration saved")
				return nil
			}
			return interactiveConfigure()
		},
	}
}

func interactiveConfigure() error {
	reader := bufio.NewReader(os.Stdin)

	output.PrintHeader("InMail CLI Configuration")
	fmt.Println()

	fmt.Print("API Key (input hidden): ")
	keyBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("read API key: %w", err)
	}
	apiKey := strings.TrimSpace(string(keyBytes))

	fmt.Printf("Base URL [https://inmail.dev/v1]: ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = "https://inmail.dev/v1"
	}

	c := &config.Config{APIKey: apiKey, BaseURL: baseURL}
	if err := config.Save(c); err != nil {
		output.PrintError(err.Error())
		return err
	}
	output.PrintSuccess("Configuration saved to ~/.inmail/config.json")
	return nil
}
