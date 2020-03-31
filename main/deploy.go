package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hekonsek/clails/clails"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
	"strings"
)
import "github.com/aws/aws-sdk-go/service/cloudformation"

func init() {
	RootCommand.AddCommand(GenerateCommand)
}

var GenerateCommand = &cobra.Command{
	Use:   "deploy",
	Short: "deploy into a cloud",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := clails.LoadProjectFromYmlFile("clails.yml")
		osexit.ExitOnError(err)
		monitoring, templates, err := clails.NewAwsDriver().Generate(project)
		osexit.ExitOnError(err)

		sess, err := newSession()
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
				fmt.Printf("CloudFormation stack for environment named %s already exists. Skipping.\n", color.GreenString(project.Name + "-monitoring"))
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
					fmt.Printf("CloudFormation stack for environment named %s already exists. Skipping.\n", color.GreenString(project.Name + "-" + env))
				} else {
					osexit.ExitOnError(err)
				}
			} else {
				fmt.Printf("CloudFormation stack %s created.\n", color.GreenString(project.Name+"-"+env))
			}
		}
	},
}

func newSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func isAlreadyExistsException(err error) bool {
	return strings.HasPrefix(err.Error(), "AlreadyExistsException")
}