package cmd

import (
	"fmt"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage domain configuration",
	}
	cmd.AddCommand(
		configGetCmd(),
		configCreateCmd(),
		configUpdateCmd(),
	)
	return cmd
}

func configGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "Get domain configuration",
		Example: `  inmail config get`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, status, err := apiClient.Get("/config")
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

func configCreateCmd() *cobra.Command {
	var attachments bool
	var catchAll bool

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create domain configuration",
		Example: `  inmail config create --attachments --no-catch-all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{
				"attachments": attachments,
				"catch_all":   catchAll,
			}
			data, status, err := apiClient.Post("/config", body)
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
	cmd.Flags().BoolVar(&attachments, "attachments", true, "Enable attachment storage")
	cmd.Flags().BoolVar(&catchAll, "catch-all", false, "Enable catch-all inbox")
	return cmd
}

func configUpdateCmd() *cobra.Command {
	var attachments bool
	var catchAll bool
	var attachmentsSet bool
	var catchAllSet bool

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update domain configuration",
		Example: `  inmail config update --attachments=false`,
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]interface{}{}
			if cmd.Flags().Changed("attachments") {
				body["attachments"] = attachments
				attachmentsSet = true
			}
			if cmd.Flags().Changed("catch-all") {
				body["catch_all"] = catchAll
				catchAllSet = true
			}
			_ = attachmentsSet
			_ = catchAllSet
			if len(body) == 0 {
				return fmt.Errorf("no flags provided; use --attachments or --catch-all")
			}
			data, status, err := apiClient.Put("/config", body)
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
	cmd.Flags().BoolVar(&attachments, "attachments", true, "Enable attachment storage")
	cmd.Flags().BoolVar(&catchAll, "catch-all", false, "Enable catch-all inbox")
	return cmd
}
