package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "1.0.0"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the inmail CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("inmail CLI v%s\n", Version)
		},
	}
}
