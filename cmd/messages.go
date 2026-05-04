package cmd

import (
	"fmt"
	"net/url"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "View and manage messages",
	}
	cmd.AddCommand(
		messagesListCmd(),
		messagesGetCmd(),
		messagesDeleteCmd(),
	)
	return cmd
}

func messagesListCmd() *cobra.Command {
	var limit int
	var cursor, account string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List messages",
		Example: `  inmail messages list
  inmail messages list --account user@example.com
  inmail messages list --limit 10 --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			if account != "" {
				q.Set("account", account)
			}
			path := "/messages"
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
	cmd.Flags().StringVar(&account, "account", "", "Filter by account email")
	return cmd
}

func messagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <id>",
		Short:   "Get a message by ID",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail messages get msg_abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, status, err := apiClient.Get("/messages/" + url.PathEscape(args[0]))
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

func messagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Short:   "Delete a message",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail messages delete msg_abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !flagNonInteractive {
				if !confirm(fmt.Sprintf("Delete message %s?", args[0])) {
					output.PrintInfo("Aborted")
					return nil
				}
			}
			data, status, err := apiClient.Delete("/messages/" + url.PathEscape(args[0]))
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
				output.PrintSuccess("Message deleted")
			}
			return nil
		},
	}
}
