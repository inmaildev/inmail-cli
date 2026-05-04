package cmd

import (
	"fmt"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newStatsCmd() *cobra.Command {
	var days int

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Get domain statistics",
		Example: `  inmail stats
  inmail stats --days 30
  inmail stats --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/stats?days=%d", days)
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
	cmd.Flags().IntVar(&days, "days", 7, "Number of days to include")
	return cmd
}
