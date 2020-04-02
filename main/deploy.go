package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/fatih/color"
	"github.com/hekonsek/aws-session"
	"github.com/hekonsek/clails/clails"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
	"strings"
)
import "github.com/aws/aws-sdk-go/service/cloudformation"

var deployCommandDryRun bool

func init() {
	deployCommand.Flags().BoolVarP(&deployCommandDryRun, "dry-run", "", false, "")
	RootCommand.AddCommand(deployCommand)
}

var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "deploy into a cloud",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := clails.LoadProjectFromYmlFile("clails.yml")
		osexit.ExitOnError(err)
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
			sess, err := aws_session.NewSession()
			if err != nil {
				osexit.ExitOnError(err)
			}
			cloudformationService := cloudformation.New(sess)

			_, err = cloudformationService.CreateStack(&cloudformation.CreateStackInput{
				StackName:    aws.String(project.Name + "-monitoring"),
				TemplateBody: aws.String(monitoring),
				Capabilities: aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
			})
			if err != nil {
				if isAlreadyExistsException(err) {
					fmt.Printf("CloudFormation stack for environment named %s already exists. Skipping.\n", color.GreenString(project.Name+"-monitoring"))
				} else {
					osexit.ExitOnError(err)
				}
			} else {
				fmt.Printf("CloudFormation stack %s created.\n", color.GreenString(project.Name+"-monitoring"))
			}

			for env, template := range templates {
				_, err = cloudformationService.CreateStack(&cloudformation.CreateStackInput{
					StackName:    aws.String(project.Name + "-" + env),
					TemplateBody: aws.String(template),
					Capabilities: aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
				})
				if err != nil {
					if isAlreadyExistsException(err) {
						fmt.Printf("CloudFormation stack for environment named %s already exists. Skipping.\n", color.GreenString(project.Name+"-"+env))
					} else {
						osexit.ExitOnError(err)
					}
				} else {
					fmt.Printf("CloudFormation stack %s created.\n", color.GreenString(project.Name+"-"+env))
				}
			}
		}
	},
}

func isAlreadyExistsException(err error) bool {
	return strings.HasPrefix(err.Error(), "AlreadyExistsException")
}