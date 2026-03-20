package cmd

import (
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

// GetRootCmd returns the root command
func GetRootCmd() (*cobra.Command, error) {
	// Higher verbosity = more log output
	var (
		verbosity int
		quiet     bool
	)

	rootCmd := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var verbosityToSet log.Level
			switch verbosity {
			case 0:
				verbosityToSet = log.WARN
			case 1:
				verbosityToSet = log.INFO
			default:
				verbosityToSet = log.DEBUG
			}
			if quiet {
				verbosityToSet = log.QUIET
			}

			return log.DefaultLogger().SetLevel(verbosityToSet)
		},
	}

	AddManPagesCmd(rootCmd)
	AddVersionCmd(rootCmd)

	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbosity", "v", "verbosity level (-v|-vv)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode")
	rootCmd.MarkFlagsMutuallyExclusive("verbosity", "quiet")
	return rootCmd, nil
}
