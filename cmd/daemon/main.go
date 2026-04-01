package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/cleanup"
	"github.com/caffeine-addictt/camserver/internal/feed"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

func main() {
	ctx, done := context.WithCancel(context.Background())
	wg := cleanup.Watch(ctx, done)
	defer func() {
		done()
		wg.Wait()
		log.Sync()
	}()

	log.SetInterruptHandler(false)
	log.DefaultLogger().SetName("camserver-daemon")

	root, err := cmd.GetRootCmd()
	if err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}

	root.Cmd.Use = "camserver-daemon"
	root.Cmd.Short = "camserver daemon"
	root.Cmd.Long = util.MultilineString(
		"Camera Server daemon",
		"",
		"Handles everything in the backend.",
	)
	root.Cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return run(wg, root.CfgManager, cmd, args)
	}

	cmd.HandleCmdExec(ctx, root.Cmd)
}

func run(wg *sync.WaitGroup, cfgManager *config.ConfigManager, c *cobra.Command, _ []string) error {
	cfgManager.RegisterCallback(func(newCfg, oldCfg *config.Config) {
		fmt.Printf("NEW\n%+v\nOLD\n%+v\n", newCfg, oldCfg)
	})

	_ = feed.NewFeedManager(c.Context(), wg)
	// defer fm.Stop()

	for {
		select {
		case <-c.Context().Done():
			return nil
		case <-time.After(time.Second * 1):
			fmt.Printf("%+v\n", cfgManager)
		}
	}
}
