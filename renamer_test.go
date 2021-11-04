package plexrenamer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
	results, err := ScanDir("artifacts", "output", true)
	assert.NoError(t, err)
	assert.Equal(t, []FileResult{
		{
			From: "artifacts/アニメ　ラブライブ！スーパースター！！（１１）「もう一度、あの場所で」 20211010.mp4",
			To:   "output/ラブライブ！スーパースター！！/Season 01/ラブライブ！スーパースター！！ S01E11 もう一度、あの場所で.mp4",
		},
	}, results)
}
