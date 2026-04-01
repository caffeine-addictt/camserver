package cmd

import (
	"fmt"

	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

type ConfigContextKeyT string

const ConfigContextKey ConfigContextKeyT = "camserver_config"

// RootCmd for convenience stuff
type RootCmd struct {
	CfgManager *config.ConfigManager
	Cmd        *cobra.Command
}

// GetRootCmd returns the root command
func GetRootCmd() (*RootCmd, error) {
	// Higher verbosity = more log output
	var (
		verbosity int
		quiet     bool

		configPath string
	)

	root := &RootCmd{}
	root.Cmd = &cobra.Command{
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

			if err := log.DefaultLogger().SetLevel(verbosityToSet); err != nil {
				return err
			}

			cm, err := config.NewConfigManager(configPath)
			if err != nil {
				return err
			}

			root.CfgManager = cm
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if root.CfgManager == nil {
				return fmt.Errorf("failed to close config manager")
			}
			root.CfgManager.Close()
			return nil
		},
	}

	AddManPagesCmd(root.Cmd)
	AddVersionCmd(root.Cmd)

	root.Cmd.PersistentFlags().CountVarP(&verbosity, "verbosity", "v", "verbosity level (-v|-vv)")
	root.Cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode")
	root.Cmd.MarkFlagsMutuallyExclusive("verbosity", "quiet")

	root.Cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file")
	return root, nil
}
