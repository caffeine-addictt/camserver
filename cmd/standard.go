package cmd

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	mango "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

func AddManPagesCmd(c *cobra.Command) {
	c.AddCommand(&cobra.Command{
		Use:                   "man",
		Short:                 "generate manpages",
		Long:                  "Generate manpages.",
		SilenceUsage:          true,
		SilenceErrors:         true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mango.NewManPage(1, c)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	})
}

func AddVersionCmd(c *cobra.Command) {
	c.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			if info, ok := debug.ReadBuildInfo(); ok {
				parts := strings.Split(info.Main.Path, "/")
				fmt.Printf("Version: %s %s\n", parts[len(parts)-1], info.Main.Version)

				for _, s := range info.Settings {
					switch s.Key {
					case "vcs.revision":
						fmt.Printf("Commit: %s\n", s.Value)
					case "vcs.modified":
						fmt.Printf("Dirty: %s\n", s.Value)
					}
				}
			}

			fmt.Printf("Go: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	})
}
