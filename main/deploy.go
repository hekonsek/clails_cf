package main

import (
	"fmt"
	"github.com/hekonsek/clails/clails"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(GenerateCommand)
}

var GenerateCommand = &cobra.Command{
	Use:   "deploy",
	Short: "deploy into a cloud",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := clails.LoadProjectFromYmlFile("clails.yml")
		osexit.ExitOnError(err)
		monitoring, templates, err := (&clails.AwsDriver{}).Generate(project)
		osexit.ExitOnError(err)

		fmt.Println(monitoring)
		fmt.Println()
		fmt.Println(templates["staging"])
	},
}
