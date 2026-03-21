package main

import (
	"time"

	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/cleanup"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

func main() {
	ctx, done, wg := cleanup.Watch()
	defer func() {
		done()
		wg.Wait()
		log.Sync()
	}()

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
	rootCmd.RunE = run

	cmd.HandleCmdExec(ctx, rootCmd)
}

func run(c *cobra.Command, args []string) error {
	cfg := config.NewConfigManager("")
	defer cfg.Close()

	select {
	case <-c.Context().Done():
	case <-time.After(time.Second * 30):
	}
	return nil
}
