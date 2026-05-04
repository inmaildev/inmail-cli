package cmd

import (
	"fmt"
	"net/url"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newSendLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-logs",
		Short: "View send logs",
	}
	cmd.AddCommand(sendLogsListCmd(), sendLogsGetCmd())
	return cmd
}

func sendLogsListCmd() *cobra.Command {
	var limit int
	var cursor string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List send logs",
		Example: `  inmail send-logs list --limit 20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("limit", fmt.Sprintf("%d", limit))
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			path := "/send-logs"
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

func sendLogsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <id>",
		Short:   "Get a send log entry by ID",
		Args:    cobra.ExactArgs(1),
		Example: `  inmail send-logs get abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, status, err := apiClient.Get("/send-logs/" + url.PathEscape(args[0]))
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
