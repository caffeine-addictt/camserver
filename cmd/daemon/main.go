package main

import (
	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/lattesec/log"
)

func main() {
	defer log.Sync()
	log.SetInterruptHandler(false)
	log.DefaultLogger().SetName("camserver-daemon")

	rootCmd, err := cmd.GetRootCmd()
	if err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}

	rootCmd.Use = "camserver-daemon"
	rootCmd.Short = "camserver daemon"
	rootCmd.Long = util.MultilineString(
		"Camera Server daemon",
		"",
		"Handles everything in the backend.",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}
}
