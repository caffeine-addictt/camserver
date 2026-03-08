package main

import (
	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func init() {
	cmd.AddManPagesCmd(rootCmd)
}

func main() {
	defer log.Sync()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}
}
