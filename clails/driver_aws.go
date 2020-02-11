package clails

import "gopkg.in/yaml.v2"

type AwsDriver struct {
}

func (*AwsDriver) Validate(*Project) error {
	return nil
}

func (driver *AwsDriver) GenerateModel(project *Project) (map[string]interface{}, error) {
	err := driver.Validate(project)
	if err != nil {
		return nil, err
	}

	templatesModels := map[string]interface{}{}
	for _, env := range []string{"staging"} {
		resources := map[string]interface{}{}
		templatesModels[env] = map[string]interface{}{
			"Resources": resources,
		}

		for _, service := range project.Services {
			if service.Type == "kafka" {
				if service.Distribution == "ami" {
					resources["KafkaServer"] = kafkaBackingServiceAmi(project)
				}
			}
		}
	}

	return templatesModels, nil
}

func kafkaBackingServiceAmi(project *Project) map[string]interface{} {
	return map[string]interface{}{
		"Type": "AWS::EC2::Instance",
		"Properties": map[string]interface{}{
			"ImageId":      "ami-0874ee9615fa7a282",
			"InstanceType": "m5.large",
			"KeyName":      "default",
			"Tags": []map[string]string{
				{"Key": "Name", "Value": project.Name + "-staging-kafka-server"},
			},
		},
	}
}

func (driver *AwsDriver) Generate(project *Project) (map[string]string, error) {
	templatesModels, err := driver.GenerateModel(project)
	if err != nil {
		return nil, err
	}

	templates := map[string]string{}
	for env, model := range templatesModels {
		envModel, err := yaml.Marshal(&model)
		if err != nil {
			return nil, err
		}
		templates[env] = string(envModel)
	}

	return templates, nil
}
