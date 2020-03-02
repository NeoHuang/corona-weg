package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJetztApi_GetCurrent(t *testing.T) {
	jetztApi := NewJetztApi(time.Millisecond)

	epidemicMap := jetztApi.GetCurrent()

	assert.Len(t, epidemicMap, 17) // 16 bundeslaender + Gesamt
}
