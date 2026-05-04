package cmd

import (
	"fmt"

	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

func newSendCmd() *cobra.Command {
	var from, subject, textBody, htmlBody string
	var to, cc []string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send an outbound email",
		Long: `Send an outbound email from your domain.

If no custom domain is configured, messages are sent from @inmail.dev.`,
		Example: `  inmail send --to user@example.com --subject "Hello" --text "Hi there"
  inmail send --to a@x.com --to b@x.com --subject "Hi" --html "<h1>Hello</h1>"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(to) == 0 {
				return fmt.Errorf("--to is required")
			}
			if subject == "" {
				return fmt.Errorf("--subject is required")
			}
			if textBody == "" && htmlBody == "" {
				return fmt.Errorf("--text or --html is required")
			}

			body := map[string]interface{}{
				"to":      to,
				"subject": subject,
			}
			if from != "" {
				body["from"] = from
			}
			if len(cc) > 0 {
				body["cc"] = cc
			}
			if textBody != "" {
				body["textBody"] = textBody
			}
			if htmlBody != "" {
				body["htmlBody"] = htmlBody
			}
			body["charset"] = "UTF-8"

			data, status, err := apiClient.Post("/send", body)
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
				output.PrintSuccess("Email sent")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "Sender address (defaults to your domain)")
	cmd.Flags().StringArrayVar(&to, "to", nil, "Recipient address (repeatable)")
	cmd.Flags().StringArrayVar(&cc, "cc", nil, "CC address (repeatable)")
	cmd.Flags().StringVar(&subject, "subject", "", "Email subject (required)")
	cmd.Flags().StringVar(&textBody, "text", "", "Plain text body")
	cmd.Flags().StringVar(&htmlBody, "html", "", "HTML body")
	return cmd
}
