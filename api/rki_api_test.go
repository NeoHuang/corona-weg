package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRkiApi_GetCurrent(t *testing.T) {
	jetztApi := NewRkiApi(time.Millisecond)

	epidemicMap, err := jetztApi.GetCurrent()
	require.NoError(t, err)

	assert.True(t, len(epidemicMap) > 0) // 16 bundeslaender + Gesamt
}
