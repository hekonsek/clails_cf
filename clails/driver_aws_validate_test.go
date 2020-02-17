package clails_test

import (
	"github.com/hekonsek/clails/clails"
	"github.com/stretchr/testify/assert"
)
import "testing"

func TestValidationShouldDefaultToAmiKafka(t *testing.T) {
	// Given
	driver := clails.NewAwsDriver()
	project, err := clails.LoadProjectFromYmlFile("../samples/kafka.yml")
	assert.NoError(t, err)

	// When
	err = driver.Validate(project)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, "ami", project.Services[0].Distribution)
}

func TestValidationShouldDefaultToStagingAndProductionEnvironments(t *testing.T) {
	// Given
	driver := clails.NewAwsDriver()
	project, err := clails.LoadProjectFromYmlFile("../samples/kafka.yml")
	assert.NoError(t, err)

	// When
	err = driver.Validate(project)
	assert.NoError(t, err)

	// Then
	assert.Len(t, project.Environments, 2)
	assert.Equal(t, []string{"staging", "production"}, project.Environments)
}
