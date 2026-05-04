package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users in your domain",
	}
	cmd.AddCommand(
		usersListCmd(),
		usersGetCmd(),
		usersCreateCmd(),
		usersUpdateCmd(),
		usersDeleteCmd(),
	)
	return cmd
}

func usersListCmd() *cobra.Command {
	var limit int
	var cursor string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Example: `  inmail users list
  inmail users list --limit 20
  inmail users list --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			path := "/users"
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
	cmd.Flags().IntVar(&limit, "limit", 50, "Number of results (max 100)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}

func usersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <email>",
		Short:   "Get a user by email",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail users get user@example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, status, err := apiClient.Get("/users/" + url.PathEscape(args[0]))
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

func usersCreateCmd() *cobra.Command {
	var email, password, name string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new user",
		Example: `  inmail users create --email user@example.com --password secret --name "John Doe"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" {
				return fmt.Errorf("--email is required")
			}
			body := map[string]string{"email": email, "password": password, "name": name}
			data, status, err := apiClient.Post("/users", body)
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
	cmd.Flags().StringVar(&email, "email", "", "User email (required)")
	cmd.Flags().StringVar(&password, "password", "", "User password")
	cmd.Flags().StringVar(&name, "name", "", "User display name")
	cmd.MarkFlagRequired("email")
	return cmd
}

func usersUpdateCmd() *cobra.Command {
	var name, password string

	cmd := &cobra.Command{
		Use:     "update <email>",
		Short:   "Update a user",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail users update user@example.com --name "Jane Doe"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]string{}
			if name != "" {
				body["name"] = name
			}
			if password != "" {
				body["password"] = password
			}
			data, status, err := apiClient.Put("/users/"+url.PathEscape(args[0]), body)
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
	cmd.Flags().StringVar(&name, "name", "", "New display name")
	cmd.Flags().StringVar(&password, "password", "", "New password")
	return cmd
}

func usersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <email>",
		Short:   "Deactivate a user",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail users delete user@example.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !flagNonInteractive {
				if !confirm(fmt.Sprintf("Deactivate user %s?", args[0])) {
					output.PrintInfo("Aborted")
					return nil
				}
			}
			data, status, err := apiClient.Delete("/users/" + url.PathEscape(args[0]))
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
				output.PrintSuccess("User deactivated")
			}
			return nil
		},
	}
}

func confirm(prompt string) bool {
	var resp string
	fmt.Printf("%s [y/N]: ", prompt)
	fmt.Scanln(&resp)
	return resp == "y" || resp == "Y"
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
