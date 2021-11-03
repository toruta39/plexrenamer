package plexrenamer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlexFormat(t *testing.T) {
	do := func(title string, season, episode int, subtitle string, date time.Time, ext, path string) {
		output, err := PlexFormat(&ProgramInfo{
			Title:    title,
			Season:   season,
			Episode:  episode,
			Subtitle: subtitle,
			Date:     date,
		}, ext)
		assert.NoError(t, err)
		assert.Equal(t, path, output)
	}

	do(
		"それＳｎｏｗ　Ｍａｎにやらせて下さい",
		1,
		-1,
		"斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！",
		time.Date(2021, time.October, 10, 0, 0, 0, 0, time.Local),
		".m2ts",
		"それＳｎｏｗ　Ｍａｎにやらせて下さい/Season 01/それＳｎｏｗ　Ｍａｎにやらせて下さい 2021-10-10 斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！..",
	)

	do(
		"マブラヴ　オルタネイティヴ",
		1,
		2,
		"",
		time.Date(2021, time.October, 14, 0, 0, 0, 0, time.Local),
		".m2ts",
		"マブラヴ　オルタネイティヴ/Season 01/マブラヴ　オルタネイティヴ S01E02.m2ts",
	)

	do(
		"ラブライブ！スーパースター！！",
		1,
		11,
		"もう一度、あの場所で",
		time.Date(2021, time.October, 10, 0, 0, 0, 0, time.Local),
		".mp4",
		"ラブライブ！スーパースター！！/Season 01/ラブライブ！スーパースター！！ S01E11 もう一度、あの場所で.mp4",
	)
}
