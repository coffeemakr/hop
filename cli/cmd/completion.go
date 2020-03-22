package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
var completionCommand = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCommand.GenBashCompletion(os.Stdout)
	},
}
