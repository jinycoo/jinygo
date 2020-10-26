package main

import (
	"log"

	"github.com/spf13/cobra"

	"jinycoo.com/jinygo/tools/jiny/commands"
)

func main() {
	cmds := &cobra.Command{
		Use:   "jiny",
		Short: "快速创建基于Jinygo框架的Golang项目，及部署配置",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	commands.AddCommands(cmds)

	if err := cmds.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}