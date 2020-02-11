package main

import (
	"fmt"
	"github.com/hekonsek/clails/clails"
	"github.com/hekonsek/clails/util"
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
		util.CliError(err)
		templates, err := (&clails.AwsDriver{}).Generate(project)
		util.CliError(err)

		fmt.Println(templates)
	},
}
