package main

import (
	"context"

	"github.com/caffeine-addictt/camserver/cmd"
	"github.com/caffeine-addictt/camserver/internal/cleanup"
	"github.com/caffeine-addictt/camserver/internal/util"
	"github.com/lattesec/log"
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
	log.DefaultLogger().SetName("camserver-web")

	root, err := cmd.GetRootCmd()
	if err != nil {
		log.Fatal().Msg(err.Error()).Send()
	}

	root.Cmd.Use = "camserver"
	root.Cmd.Short = "camserver web"
	root.Cmd.Long = util.MultilineString(
		"Camera Server Web Interface",
		"",
		"Access bridge to camserver-daemon",
	)

	cmd.HandleCmdExec(ctx, root.Cmd)
}
