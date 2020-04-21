package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/clails/clails"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var deployCommandDryRun bool

func init() {
	deployCommand.Flags().BoolVarP(&deployCommandDryRun, "dry-run", "", false, "")
	RootCommand.AddCommand(deployCommand)
}

var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "deploy Cloud Formation stack into AWS",
	Run: func(cmd *cobra.Command, args []string) {
		ymlProject := clails.NewYmlProject()
		project, err := ymlProject.LoadFromFile("clails.yml")
		osexit.ExitOnError(ymlProject.FriendlyMessage(err))

		monitoring, templates, err := clails.NewAwsDriver().Generate(project)
		osexit.ExitOnError(err)

		if deployCommandDryRun {
			fmt.Printf("CloudFormation stack for %s environment:\n", color.GreenString("monitoring"))
			fmt.Println()
			fmt.Println(monitoring)
			fmt.Println("---")
			for env, template := range templates {
				fmt.Printf("CloudFormation stack for %s environment:\n", color.GreenString(env))
				fmt.Println()
				fmt.Println(template)
				fmt.Println()
				fmt.Println("---")
				fmt.Println()
			}
		} else {
			var deploymentNotifications = make(chan clails.EnvironmentStatus)
			go func() {
				err = clails.NewDeployer().Deploy(project, monitoring, templates, deploymentNotifications)
				if err != nil {
					osexit.ExitOnError(err)
				}
			}()
			for notification := range deploymentNotifications {
				stackName := fmt.Sprintf("%s-%s", project.Name, notification.Environment)
				if notification.Succedded {
					fmt.Printf("CloudFormation stack %s created.\n", color.GreenString(stackName))
				} else {
					fmt.Printf("CloudFormation stack %s already exists. Skipping.\n", color.GreenString(stackName))
				}
			}
		}
	},
}
