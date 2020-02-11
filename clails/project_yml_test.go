package clails_test

import (
	"github.com/hekonsek/clails/clails"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadingNoSuchFile(t *testing.T) {
	// When
	_, err := clails.LoadProjectFromYmlFile("noSuchFile.yml")

	// Then
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "no such file: noSuchFile.yml"))
}

func TestLoadingServices(t *testing.T) {
	// When
	project, err := clails.LoadProjectFromYmlFile("../samples/kafka-ami.yml")
	assert.NoError(t, err)

	// Then
	assert.NotNil(t, project.Services)
	assert.Len(t, project.Services, 1)
}
