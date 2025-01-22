package root

import (
	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd"
	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd/restful"
	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd/version"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cmd.Factory) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "listen-tube <command> <subcommand> [flags]",
		Short: "listen-tube",
		Long:  `listen-tube is a tool collection for bash environments.`,

		Annotations: map[string]string{
			"version": "0.0.1",
			"website": "www.github.com/gogodjzhu/listen-tube.xyz",
		},
	}

	cmd.AddCommand(version.NewCmdVersion(f))

	restfulCmd := restful.NewCmdRestful(f)
	cmd.AddCommand(restfulCmd)

	return cmd, nil
}
