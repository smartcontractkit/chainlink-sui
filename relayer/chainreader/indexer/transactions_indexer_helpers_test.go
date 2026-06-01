package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExecuteInitFunction(t *testing.T) {
	assert.True(t, isExecuteInitFunction(executeFunctionV1))
	assert.True(t, isExecuteInitFunction(executeFunctionV2))
	assert.False(t, isExecuteInitFunction("finish_execute"))
	assert.False(t, isExecuteInitFunction("finish_execute_v2"))
}
