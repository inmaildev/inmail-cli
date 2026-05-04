package cmd

import (
	"fmt"
	"net/url"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manage email accounts",
	}
	cmd.AddCommand(
		accountsListCmd(),
		accountsGetCmd(),
		accountsCreateCmd(),
		accountsUpdateCmd(),
		accountsDeleteCmd(),
	)
	return cmd
}

func accountsListCmd() *cobra.Command {
	var limit int
	var cursor string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all accounts",
		Example: `  inmail accounts list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			path := "/accounts"
			if len(q) > 0 {
				path += "?" + q.Encode()
			}
			data, status, err := apiClient.Get(path)
			if err != nil {
				output.PrintError(err.Error())
				return err
			}
			if status >= 400 {
				output.PrintAPIError(status, data)
				return fmt.Errorf("API error %d", status)
			}
			output.PrintJSON(data)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 50, "Number of results")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}

func accountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <id>",
		Short:   "Get an account by ID",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail accounts get acc_123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, status, err := apiClient.Get("/accounts/" + url.PathEscape(args[0]))
			if err != nil {
				output.PrintError(err.Error())
				return err
			}
			if status >= 400 {
				output.PrintAPIError(status, data)
				return fmt.Errorf("API error %d", status)
			}
			output.PrintJSON(data)
			return nil
		},
	}
}

func accountsCreateCmd() *cobra.Command {
	var email string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new account",
		Example: `  inmail accounts create --email inbox@yourdomain.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]string{"email": email}
			data, status, err := apiClient.Post("/accounts", body)
			if err != nil {
				output.PrintError(err.Error())
				return err
			}
			if status >= 400 {
				output.PrintAPIError(status, data)
				return fmt.Errorf("API error %d", status)
			}
			output.PrintJSON(data)
			return nil
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "Account email (required)")
	cmd.MarkFlagRequired("email")
	return cmd
}

func accountsUpdateCmd() *cobra.Command {
	var email string

	cmd := &cobra.Command{
		Use:     "update <id>",
		Short:   "Update an account",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail accounts update acc_123 --email new@yourdomain.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]string{"email": email}
			data, status, err := apiClient.Put("/accounts/"+url.PathEscape(args[0]), body)
			if err != nil {
				output.PrintError(err.Error())
				return err
			}
			if status >= 400 {
				output.PrintAPIError(status, data)
				return fmt.Errorf("API error %d", status)
			}
			output.PrintJSON(data)
			return nil
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "New email address")
	return cmd
}

func accountsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Short:   "Delete an account",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail accounts delete acc_123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !flagNonInteractive {
				if !confirm(fmt.Sprintf("Delete account %s?", args[0])) {
					output.PrintInfo("Aborted")
					return nil
				}
			}
			data, status, err := apiClient.Delete("/accounts/" + url.PathEscape(args[0]))
			if err != nil {
				output.PrintError(err.Error())
				return err
			}
			if status >= 400 {
				output.PrintAPIError(status, data)
				return fmt.Errorf("API error %d", status)
			}
			if len(data) > 0 {
				output.PrintJSON(data)
			} else {
				output.PrintSuccess("Account deleted")
			}
			return nil
		},
	}
}
