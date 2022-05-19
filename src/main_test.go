package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunApplication(t *testing.T) {
	exitCode := runApplication()
	assert.Equal(t, 0, exitCode)
}
