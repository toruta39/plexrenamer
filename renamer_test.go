package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
	toFiles, err := ScanDir("artifacts", "output", true)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"output/アニメ\u3000ラブライブ！スーパースター！！/Season 01/アニメ\u3000ラブライブ！スーパースター！！ S01E11 もう一度、あの場所で.mp4",
	}, toFiles)
}
