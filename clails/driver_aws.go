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
				return errors.New("unknown Kafka service distribution: " + service.Distribution)
			}
		} else if service.Type == "kubernetes" {
			if service.Distribution == "" {
				service.Distribution = "eks-nodegroup"
			}
			if service.Distribution != "eks-nodegroup" {
				return errors.New("unknown kubernetes service distribution: " + service.Distribution)
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
			"Parameters" : map[string]map[string]interface{}{
				"DefaultVpcSubnetsIds" : {
					"Type" : "CommaDelimitedList",
					"Default" : "",
				},
			},
			"Resources": resources,
		}

		for _, service := range project.Services {
			if service.Type == "kafka" {
				if service.Distribution == "ami" {
					resources["KafkaServer"] = kafkaBackingServiceAmi(project, env)
				}
			} else if service.Type == "kubernetes" {
				if service.Distribution == "eks-nodegroup" {
					for k, v := range eksClusterBackingServiceAmi(project, env) {
						resources[k] = v
					}
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

func eksClusterBackingServiceAmi(project *Project, env string) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"EksCluster": {
			"Type": "AWS::EKS::Cluster",
			"Properties": map[string]interface{}{
				"Name":    project.Name + "-" + env,
				"RoleArn": map[string][]string{
					"Fn::GetAtt":{ "EksClusterRole", "Arn" },
				},
				"ResourcesVpcConfig": map[string]interface{}{
					"SubnetIds": map[string]string{
						"Ref": "DefaultVpcSubnetsIds",
					},
				},
			},
		},
		"EksClusterRole": {
			"Type": "AWS::IAM::Role",
			"Properties": map[string]interface{}{
				"RoleName": project.Name + "-EksCluster-" + "-" + env,
				"AssumeRolePolicyDocument": map[string]interface{}{
					"Statement": []interface{}{
						map[string]interface{}{
							"Sid":    "",
							"Effect": "Allow",
							"Principal": map[string]interface{}{
								"Service": "eks.amazonaws.com",
							},
							"Action": []string{"sts:AssumeRole"},
						},
					},
					"Version": "2012-10-17",
				},
				"ManagedPolicyArns": []string{
					"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
					"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
				},
			},
		},
		"EksNodeGroup": {
			"Type": "AWS::EKS::Nodegroup",
			"Properties": map[string]interface{}{
				"ClusterName": project.Name + "-" + env,
				"NodeRole": map[string][]string{
					"Fn::GetAtt":{ "EksNodeGroupRole", "Arn" },
				},
				"Subnets": map[string]string{
					"Ref": "DefaultVpcSubnetsIds",
				},
			},
			"DependsOn": "EksCluster",
		},
		"EksNodeGroupRole": {
			"Type": "AWS::IAM::Role",
			"Properties": map[string]interface{}{
				"RoleName": project.Name + "-EksNodeGroup" + "-" + env,
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
					"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
					"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
					"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
				},
			},
		},
	}
}

func monitoringTemplate(project *Project) map[string]map[string]interface{} {
	roleName := project.Name + "-monitoring-ec2-read-access"
	profileName := project.Name + "-monitoring-ec2"
	return map[string]map[string]interface{}{
		"MonitoringSecurityGroup": {
			"Type": "AWS::EC2::SecurityGroup",
			"Properties": map[string]interface{}{
				"GroupName":        project.Name + "-monitoring",
				"GroupDescription": project.Name + "-monitoring",
				"SecurityGroupIngress": []map[string]interface{}{
					{
						"IpProtocol": "tcp",
						"FromPort":   9090,
						"ToPort":     9090,
						"CidrIp":     "0.0.0.0/0",
					},
					{
						"IpProtocol": "tcp",
						"FromPort":   3000,
						"ToPort":     3000,
						"CidrIp":     "0.0.0.0/0",
					},
				},
			},
		},
		"MonitoringServer": {
			"Type": "AWS::EC2::Instance",
			"Properties": map[string]interface{}{
				"ImageId":      "ami-07a0a263711b54ac0",
				"InstanceType": "m5.large",
				"KeyName":      "default",
				"SecurityGroupIds": []interface{}{
					map[string]string{
						"Ref": "MonitoringSecurityGroup",
					},
				},
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
