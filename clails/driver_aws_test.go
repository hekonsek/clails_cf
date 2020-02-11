package clails_test

import (
	"github.com/hekonsek/clails/clails"
	"github.com/stretchr/testify/assert"
)
import "testing"

func TestGenerateStagingModel(t *testing.T) {
	// Given
	driver := &clails.AwsDriver{}
	project, err := clails.LoadProjectFromYmlFile("../samples/kafka-ami.yml")
	assert.NoError(t, err)

	// When
	model, err := driver.GenerateModel(project)
	assert.NoError(t, err)

	// Then
	assert.NotNil(t, model["staging"])
}
