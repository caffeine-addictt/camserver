package main

import (
	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/lattesec/log"
)

func main() {
	defer log.Sync()
	log.SetInterruptHandler(false)
	log.DefaultLogger().SetName("camserver-web")

	rootCmd, err := cmd.GetRootCmd()
	if err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}

	rootCmd.Use = "camserver"
	rootCmd.Short = "camserver web"
	rootCmd.Long = util.MultilineString(
		"Camera Server Web Interface",
		"",
		"Access bridge to camserver-daemon",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}
}
