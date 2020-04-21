package clails

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	aws_session "github.com/hekonsek/aws-session"
	"strings"
)

type Deployer struct {

}

func NewDeployer() *Deployer {
	return &Deployer{}
}

type EnvironmentStatus struct {
	Environment string
	Succedded bool
}

func (deployer *Deployer) Deploy(project *Project, monitoring string, templates map[string]string, envStatus chan EnvironmentStatus) error {
	sess, err := aws_session.NewSession()
	if err != nil {
		return err
	}
	cloudformationService := cloudformation.New(sess)

	_, err = cloudformationService.CreateStack(&cloudformation.CreateStackInput{
		StackName:    aws.String(project.Name + "-monitoring"),
		TemplateBody: aws.String(monitoring),
		Capabilities: aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
	})
	if err != nil {
		if isAlreadyExistsException(err) {
			envStatus <- EnvironmentStatus{Environment: "monitoring", Succedded: false}
		} else {
			return err
		}
	} else {
		envStatus <- EnvironmentStatus{Environment: "monitoring", Succedded: true}
	}

	ec2Service := ec2.New(sess)
	vpcsResults, err := ec2Service.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return err
	}
	defaultVpcId := ""
	for _, vpc := range vpcsResults.Vpcs {
		if *vpc.IsDefault == true {
			defaultVpcId = *vpc.VpcId
		}
	}
	defaultSubnetsResults, err := ec2Service.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: aws.StringSlice([]string{defaultVpcId}),
			},
		},
	})
	if err != nil {
		return err
	}
	var defaultSubnetsIds []string
	for _, subnet := range defaultSubnetsResults.Subnets {
		defaultSubnetsIds = append(defaultSubnetsIds, *subnet.SubnetId)
	}

	for env, template := range templates {
		_, err = cloudformationService.CreateStack(&cloudformation.CreateStackInput{
			StackName:    aws.String(project.Name + "-" + env),
			TemplateBody: aws.String(template),
			Capabilities: aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
			Parameters: []*cloudformation.Parameter{
				{
					ParameterKey: aws.String("DefaultVpcSubnetsIds"),
					ParameterValue: aws.String(strings.Join(defaultSubnetsIds, ",")),
				},
			},
		})
		if err != nil {
			if isAlreadyExistsException(err) {
				envStatus <- EnvironmentStatus{Environment: env, Succedded: false}
			} else {
				return err
			}
		} else {
			envStatus <- EnvironmentStatus{Environment: env, Succedded: true}
		}
	}

	close(envStatus)
	return nil
}

func isAlreadyExistsException(err error) bool {
	return strings.HasPrefix(err.Error(), "AlreadyExistsException")
}