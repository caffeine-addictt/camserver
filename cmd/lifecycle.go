package cmd

import (
	"context"

	"github.com/caffeine-addictt/camserver/internal/cleanup"
	"github.com/lattesec/log"
	"github.com/spf13/cobra"
)

func ExecCmdContext(ctx context.Context, c *cobra.Command) (*cobra.Command, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msgf("panic: %v", err).Send()
		}
	}()

	return c.ExecuteContextC(ctx)
}

func HandleCmdExec(ctx context.Context, c *cobra.Command) {
	_ = HandleCmdExecE(ctx, c)
}

func HandleCmdExecE(ctx context.Context, c *cobra.Command) error {
	_, err := ExecCmdContext(ctx, c)
	cleanup.Cleanup()
	if err != nil {
		cleanup.CleanupError()
		log.Error().Msg(err.Error()).Send()
	}
	return err
}
