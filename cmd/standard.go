package cmd

import (
	"fmt"

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
