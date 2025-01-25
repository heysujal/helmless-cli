package helmless

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helmless",
		Short: "Helmless is a tool for managing Helm charts",
	}

	cmd.AddCommand(newCreateCmd())

	return cmd
}

func Execute() error {
	return NewRootCmd().Execute()
}
