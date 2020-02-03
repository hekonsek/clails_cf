package main

import (
	"github.com/hekonsek/clails/util"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "clails",
	Short: "Clails - cloud made easy",

	Run: func(cmd *cobra.Command, args []string) {
		util.ExitOnCliError(cmd.Help())
	},
}

func main() {
	util.ExitOnCliError(RootCommand.Execute())
}