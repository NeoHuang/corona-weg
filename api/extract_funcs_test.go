package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRkiApi_GetCurrent(t *testing.T) {
	jetztApi := NewRkiApi(time.Millisecond)

	epidemicMap, err := jetztApi.GetCurrent()
	require.NoError(t, err)
	fmt.Println(epidemicMap)

	// 16 bundeslaender + Gesamt. missing if 0 infections
	assert.True(t, len(epidemicMap) > 15)
	assert.True(t, len(epidemicMap) < 18)
}

func TestJetztApi_GetCurrent(t *testing.T) {
	jetztApi := NewJetztApi(time.Millisecond)

	epidemicMap, err := jetztApi.GetCurrent()
	require.NoError(t, err)

	assert.Len(t, epidemicMap, 16) // 16 bundeslaender
}
