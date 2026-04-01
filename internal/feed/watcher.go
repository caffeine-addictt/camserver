package feed

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/lattesec/log"
)

const (
	rotationSeconds = 60 // rotate every 60s
)

type Watcher struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	ctxDone context.CancelFunc

	Camera     *config.CameraCfg
	ArchiveDir string
}

func NewWatcher(ctx context.Context, rootWg *sync.WaitGroup, cam *config.CameraCfg, archiveDir string) *Watcher {
	wCtx, wCtxDone := context.WithCancel(context.Background())
	w := &Watcher{
		wg:      &sync.WaitGroup{},
		ctx:     wCtx,
		ctxDone: wCtxDone,

		Camera:     cam,
		ArchiveDir: archiveDir,
	}

	rootWg.Go(func() {
		<-ctx.Done()
		w.Stop()
	})

	return w
}

// Stop signals the watcher to stop
func (w *Watcher) Stop() {
	w.ctxDone()
	w.wg.Wait()
}

func (w *Watcher) L() *log.LogMessage {
	return log.Info().
		WithMeta("feed", "watcher").
		WithMeta("cam", w.Camera.Name)
}

func (w *Watcher) Start() {
	w.wg.Go(func() {
		for {
			select {
			case <-w.ctx.Done():
				w.L().Info().Msg("shutting down").Send()
				return
			default:
				w.runFFmpegSession()
			}
		}
	})
}

// runFFmpegSession runs ffmpeg for a single MP4 segment and handles rotation
func (w *Watcher) runFFmpegSession() {
	timestamp := time.Now().Format("20060102-150405")
	rootDir := filepath.Join(w.ArchiveDir, w.Camera.GetDirRel())
	tmpPath := filepath.Join(rootDir, fmt.Sprintf("%s.tmp.mp4", timestamp))
	finalPath := filepath.Join(rootDir, fmt.Sprintf("%s.mp4", timestamp))

	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		w.L().Error().Msgf("failed to create archive dir at %s: %v", w.ArchiveDir, err).Send()
		return
	}
	w.L().Debug().Msgf("writing to temp %s", tmpPath).Send()

	args := []string{
		"-rtsp_transport", "tcp",
		"-i", w.Camera.Rtsp.String(),
		"-c:v", "copy",
		"-c:a", "aac",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"-t", fmt.Sprintf("%d", rotationSeconds), // duration per file
		tmpPath,
	}

	cmd := exec.CommandContext(w.ctx, "ffmpeg", args...)
	w.L().Debug().Msgf("starting ffmpeg with [%v]", args).Send()
	w.L().Info().Msg("starting ffmpeg").Send()

	if err := cmd.Run(); err != nil {
		w.L().Error().Msgf("ffmpeg exited with error: %v", err).Send()
		return
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		w.L().Error().Msgf("failed to promote tmp at %s → %s: %v", tmpPath, finalPath, err).Send()
		return
	}

	w.L().Debug().Msgf("saved segment: %s", finalPath).Send()

	// TODO: cleanup older files
}
