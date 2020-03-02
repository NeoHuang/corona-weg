package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJetztApi_GetCurrent(t *testing.T) {
	jetztApi := NewJetztApi(time.Millisecond)

	epidemicMap, err := jetztApi.GetCurrent()
	require.NoError(t, err)

	assert.Len(t, epidemicMap, 17) // 16 bundeslaender + Gesamt
}
