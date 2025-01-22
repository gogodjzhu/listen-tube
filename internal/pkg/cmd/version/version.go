package version

import (
	"fmt"

	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *cmd.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(f.IOStreams.Out, cmd.Root().Annotations["version"])
		},
	}

	return cmd
}
