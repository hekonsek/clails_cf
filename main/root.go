package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "clails",
	Short: "Clails - cloud made easy",

	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}

func main() {
	osexit.ExitOnError(RootCommand.Execute())
}
