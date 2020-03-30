package clails

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

type AwsDriver struct {
}

func NewAwsDriver() *AwsDriver {
	return &AwsDriver{}
}

var defaultEnvironments = []string{"staging", "production"}

func (*AwsDriver) Validate(project *Project) error {
	if project.Environments == nil || len(project.Environments) == 0 {
		project.Environments = defaultEnvironments
	}

	for _, service := range project.Services {
		if service.Type == "kafka" {
			if service.Distribution == "" {
				service.Distribution = "ami"
			}
			if service.Distribution != "ami" {
				return errors.New("unknown Kafka service type: " + service.Distribution)
			}
		}
	}

	return nil
}

func (driver *AwsDriver) GenerateModel(project *Project) (monitoring map[string]interface{}, environments map[string]interface{}, err error) {
	err = driver.Validate(project)
	if err != nil {
		return nil, nil, err
	}

	templatesModels := map[string]interface{}{}
	for _, env := range project.Environments {
		resources := map[string]interface{}{}
		templatesModels[env] = map[string]interface{}{
			"Resources": resources,
		}

		for _, service := range project.Services {
			if service.Type == "kafka" {
				if service.Distribution == "ami" {
					resources["KafkaServer"] = kafkaBackingServiceAmi(project, env)
				}
			}
		}
	}

	monitoringModel := map[string]interface{}{
		"Resources": monitoringTemplate(project),
	}

	return monitoringModel, templatesModels, nil
}

func (driver *AwsDriver) Generate(project *Project) (monitoring string, environments map[string]string, err error) {
	monitoringModel, templatesModels, err := driver.GenerateModel(project)
	if err != nil {
		return "", nil, err
	}

	templates := map[string]string{}
	for env, model := range templatesModels {
		envModel, err := yaml.Marshal(&model)
		if err != nil {
			return "", nil, err
		}
		templates[env] = string(envModel)
	}

	monitoringTemplate, err := yaml.Marshal(&monitoringModel)
	if err != nil {
		return "", nil, err
	}

	return string(monitoringTemplate), templates, nil
}

// Model generation

func kafkaBackingServiceAmi(project *Project, env string) map[string]interface{} {
	return map[string]interface{}{
		"Type": "AWS::EC2::Instance",
		"Properties": map[string]interface{}{
			"ImageId":      "ami-0874ee9615fa7a282",
			"InstanceType": "m5.large",
			"KeyName":      "default",
			"Tags": []map[string]string{
				{"Key": "Name", "Value": fmt.Sprintf("%s-%s-kafka-server", project.Name, env)},
			},
		},
	}
}

func monitoringTemplate(project *Project) map[string]map[string]interface{} {
	roleName := project.Name + "-monitoring-ec2-read-access"
	profileName := project.Name + "-monitoring-ec2"
	return map[string]map[string]interface{}{
		"MonitoringServer": {
			"Type": "AWS::EC2::Instance",
			"Properties": map[string]interface{}{
				"ImageId":      "ami-07a0a263711b54ac0",
				"InstanceType": "m5.large",
				"KeyName":      "default",
				"Tags": []map[string]string{
					{"Key": "Name", "Value": project.Name + "-monitoring"},
				},
				"IamInstanceProfile": profileName,
			},
		},
		"MonitoringEc2ReadAccessRole": {
			"Type": "AWS::IAM::Role",
			"Properties": map[string]interface{}{
				"RoleName": roleName,
				"AssumeRolePolicyDocument": map[string]interface{}{
					"Statement": []interface{}{
						map[string]interface{}{
							"Sid":    "",
							"Effect": "Allow",
							"Principal": map[string]interface{}{
								"Service": "ec2.amazonaws.com",
							},
							"Action": []string{"sts:AssumeRole"},
						},
					},
					"Version": "2012-10-17",
				},
				"ManagedPolicyArns": []string{
					"arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess",
				},
			},
		},
		"MonitoringEc2InstanceProfile": {
			"Type": "AWS::IAM::InstanceProfile",
			"Properties": map[string]interface{}{
				"InstanceProfileName": profileName,
				"Roles":               []string{roleName},
			},
		},
	}
}
