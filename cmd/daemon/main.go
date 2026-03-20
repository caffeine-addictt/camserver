package main

import (
	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

// Higher verbosity = more log output
var (
	verbosity int
	quiet     bool
)

var rootCmd = &cobra.Command{
	Use:           "camserver-daemon",
	Short:         "camserver daemon",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: util.MultilineString(
		"Camera Server daemon",
		"",
		"Handles everything in the backend.",
	),
}

func init() {
	cmd.AddManPagesCmd(rootCmd)
	cmd.AddVersionCmd(rootCmd)

	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbosity", "v", "verbosity level (-v|-vv)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode")
	rootCmd.MarkFlagsMutuallyExclusive("verbosity", "quiet")
}

func main() {
	defer log.Sync()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}
}
