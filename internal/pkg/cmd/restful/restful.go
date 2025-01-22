package restful

import (
	"fmt"

	"github.com/gogodjzhu/listen-tube/internal/pkg/cmd"
	"github.com/gogodjzhu/listen-tube/internal/pkg/conf"
	"github.com/gogodjzhu/listen-tube/web/controller"
	"github.com/spf13/cobra"
)

func NewCmdRestful(f *cmd.Factory) *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:    "restful",
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			if configPath == "" {
				fmt.Fprintln(f.IOStreams.Out, "[Err] config path is required")
				return
			}
			fmt.Fprintln(f.IOStreams.Out, "[Info] restful start with config path", configPath)
			config, err := conf.ReadConfig(configPath)
			if err != nil {
				fmt.Fprintln(f.IOStreams.Out, "[Err] read config failed", err)
				return
			}
			ctx := cmd.Context()
			c, err := controller.NewController(ctx, config)
			if err != nil {
				fmt.Fprintln(f.IOStreams.Out, "[Err] create controller failed", err)
				return
			}
			c.Start()
		},
	}
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")

	return cmd
}
