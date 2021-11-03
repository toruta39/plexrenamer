package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFlags(t *testing.T) {
	assert.Equal(t, "王様ランキング＜ノイタミナ＞　＃０２| 20211022",
		sanitizeFlags("王様ランキング＜ノイタミナ＞　＃０２[字] 20211022"))
}

func TestSanitizeSlots(t *testing.T) {
	assert.Equal(t, "王様ランキング|　＃０２[字] 20211022",
		sanitizeSlots("王様ランキング＜ノイタミナ＞　＃０２[字] 20211022"))
}

func TestNormalizeTitle(t *testing.T) {
	assert.Equal(t, "劇的に沈黙",
		normalizeTitle("ドラマ『劇的に沈黙』"))
}

func TestNormalizeNumber(t *testing.T) {
	assert.Equal(t, "0123456789",
		normalizeNumber("０１２３４５６７８９"))
}

func TestCommonDelimit(t *testing.T) {
	assert.Equal(t, "それＳｎｏｗ　Ｍａｎにやらせて下さい|斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！",
		commonDelimit("それＳｎｏｗ　Ｍａｎにやらせて下さい★斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！"))
}

func TestDedupeDelimiter(t *testing.T) {
	assert.Equal(t, "a|b|c",
		dedupeDelimiter("||a|b | |c|| "))
}

func TestMatchDate(t *testing.T) {
	date, output, err := matchDate("月とライカと吸血姫　第２話「宇宙飛行士への道」 20211011")
	assert.NoError(t, err)
	d, _ := time.Parse("2006-01-02", "2021-10-11")
	assert.Equal(t, d, date)
	assert.Equal(t, "月とライカと吸血姫　第２話「宇宙飛行士への道」", output)
}

func TestMatchEpisode(t *testing.T) {
	episode, output, err := matchEpisode("月とライカと吸血姫　第２話「宇宙飛行士への道」 20211011")
	assert.NoError(t, err)
	assert.Equal(t, 2, episode)
	assert.Equal(t, "月とライカと吸血姫|「宇宙飛行士への道」 20211011", output)

	episode, output, err = matchEpisode("ＶＳ魂【スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動】[字] 20211007")
	assert.NoError(t, err)
	assert.Equal(t, -1, episode)
	assert.Equal(t, "ＶＳ魂【スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動】[字] 20211007", output)
}

func TestMatchSubtitle(t *testing.T) {
	subtitle, output, err := matchSubtitle("月とライカと吸血姫|「宇宙飛行士への道」")
	assert.NoError(t, err)
	assert.Equal(t, "宇宙飛行士への道", subtitle)
	assert.Equal(t, "月とライカと吸血姫", output)

	subtitle, output, err = matchSubtitle("ＶＳ魂【スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動】")
	assert.NoError(t, err)
	assert.Equal(t, "スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動", subtitle)
	assert.Equal(t, "ＶＳ魂", output)
}

func TestMatchTitle(t *testing.T) {
	title, subtitle, output, err := matchTitle("月とライカと吸血姫")
	assert.NoError(t, err)
	assert.Equal(t, "月とライカと吸血姫", title)
	assert.Equal(t, "", subtitle)
	assert.Equal(t, "", output)

	title, subtitle, output, err = matchTitle("冒険少年 キャンプ飯対決に女優石田ひかり・志田彩良＆ロバート馬場初参戦VS池崎")
	assert.NoError(t, err)
	assert.Equal(t, "冒険少年", title)
	assert.Equal(t, "キャンプ飯対決に女優石田ひかり・志田彩良＆ロバート馬場初参戦VS池崎", subtitle)
	assert.Equal(t, "", output)

	title, subtitle, output, err = matchTitle("♪オトラクションＳＰ♪新企画「厨房の音で料理を当てろ」卓球張本＆SixTONES")
	assert.NoError(t, err)
	assert.Equal(t, "♪オトラクションＳＰ♪新企画「厨房の音で料理を当てろ」卓球張本＆SixTONES", title)
	assert.Equal(t, "", subtitle)
	assert.Equal(t, "", output)
}

func TestParse(t *testing.T) {
	t.Skip()

	do := func(input, title string, season, episode int, subtitle, date string) {
		info, err := Parse(input)
		assert.NoError(t, err)

		d, _ := time.Parse("2006-01-02", date)
		assert.Equal(t, ProgramInfo{
			Title:    title,
			Season:   season,
			Episode:  episode,
			Subtitle: subtitle,
			Date:     d,
		}, *info)
	}

	// anime
	do(
		"月とライカと吸血姫　第２話「宇宙飛行士への道」 20211011",
		"月とライカと吸血姫",
		1,
		2,
		"宇宙飛行士への道",
		"2021-10-11",
	)

	do(
		"[新]さんかく窓の外側は夜　＃１「出逢」 20211003",
		"さんかく窓の外側は夜",
		1,
		1,
		"出逢",
		"2021-10-03",
	)

	do(
		"[新]無職転生　～異世界行ったら本気だす～　＃１２「魔眼を持つ女」 20211004",
		"無職転生　～異世界行ったら本気だす～",
		1,
		12,
		"魔眼を持つ女",
		"2021-10-04",
	)

	do(
		"王様ランキング＜ノイタミナ＞　＃０２[字] 20211022",
		"王様ランキング",
		1,
		2,
		"",
		"2021-10-22",
	)

	do(
		"マブラヴ　オルタネイティヴ＜＋Ｕｌｔｒａ＞　＃０２ 20211014",
		"マブラヴ　オルタネイティヴ",
		1,
		2,
		"",
		"2021-10-14",
	)

	do(
		"アニメ　ラブライブ！スーパースター！！（１１）「もう一度、あの場所で」 20211010",
		"ラブライブ！スーパースター！！",
		1,
		11,
		"もう一度、あの場所で",
		"2021-10-10",
	)

	// tv drama
	do(
		"消えた初恋　＃３　好きな人の自宅で勉強会!？[字] 20211023",
		"消えた初恋",
		1,
		3,
		"好きな人の自宅で勉強会!？",
		"2021-10-23",
	)

	do(
		"[新]ドラマシリーズ　お耳に合いましたら。[再]　第１話　主演：伊藤万理華 20210715",
		"お耳に合いましたら。",
		1,
		1,
		"主演：伊藤万理華",
		"2021-07-15",
	)

	do(
		"[新]ドラマ『劇的に沈黙』　★第１話 20210706",
		"劇的に沈黙",
		1,
		1,
		"",
		"2021-07-06",
	)

	do(
		"ドラマ２５　ソロ活女子のススメ　第11話「ソロバーベキュー」[字] 20210612",
		"ソロ活女子のススメ",
		1,
		11,
		"ソロバーベキュー",
		"2021-06-12",
	)

	// variety
	do(
		"人志松本の酒のツマミになる話【酒井美紀＆井上咲楽の心に響いた歌詞に松本感激！】[字] 20210716",
		"人志松本の酒のツマミになる話",
		1,
		-1,
		"酒井美紀＆井上咲楽の心に響いた歌詞に松本感激！",
		"2021-07-16",
	)

	do(
		"１億３０００万人のＳＨＯＷチャンネル【名店レシピ！中華風カレーの完コピなるか】[字] 20211016",
		"１億３０００万人のＳＨＯＷチャンネル",
		1,
		-1,
		"名店レシピ！中華風カレーの完コピなるか",
		"2021-10-16",
	)

	do(
		"冒険少年 キャンプ飯対決に女優石田ひかり・志田彩良＆ロバート馬場初参戦VS池崎[字] 20211018",
		"冒険少年",
		1,
		-1,
		"キャンプ飯対決に女優石田ひかり・志田彩良＆ロバート馬場初参戦VS池崎",
		"2021-10-18",
	)

	do(
		"ＶＳ魂【スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動】[字] 20211007",
		"ＶＳ魂",
		1,
		-1,
		"スシローでオーダー対決＆岸ｖｓＩＫＫＯ料理対決▽相葉発案の新企画始動",
		"2021-10-07",
	)

	do(
		"それＳｎｏｗ　Ｍａｎにやらせて下さい★斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！ 20211010",
		"それＳｎｏｗ　Ｍａｎにやらせて下さい",
		1,
		-1,
		"斎藤工が大絶賛！最強の海鮮丼かけ頂上決戦！",
		"2021-10-10",
	)

	do(
		"オオカミ少年☆ドラゴン桜＆ももクロ参戦！あばれるドッキリ連発☆板野友美は愛犬家？ 20210611",
		"オオカミ少年",
		1,
		-1,
		"ドラゴン桜＆ももクロ参戦！あばれるドッキリ連発",
		"2021-06-11",
	)

	// news
	do(
		"ラヴィット！【ジャニーズWEST桐山＆Snow Man佐久間スタジオ登場！】[字] 20211005",
		"ラヴィット！",
		1,
		-1,
		"ジャニーズWEST桐山＆Snow Man佐久間スタジオ登場！",
		"2021-10-05",
	)

	do(
		"ｎｅｗｓ　ｚｅｒｏ[字] 東京感染４９人ことし「最少」に…去年６月以来▽櫻井翔 20211011",
		"ｎｅｗｓ　ｚｅｒｏ",
		1,
		-1,
		"東京感染４９人ことし「最少」に…去年６月以来▽櫻井翔",
		"2021-10-11",
	)

	do(
		"ＺＩＰ！[デ] 20210927",
		"ＺＩＰ！",
		1,
		-1,
		"",
		"2021-09-27",
	)

	// sp
	do(
		"２０２１ＦＮＳ歌謡祭　夏【司会は相葉雅紀！今届けたい想いを音楽にのせて】[字][デ] 20210714",
		"２０２１ＦＮＳ歌謡祭　夏",
		1,
		-1,
		"司会は相葉雅紀！今届けたい想いを音楽にのせて",
		"2021-07-14",
	)

	do(
		"痛快ＴＶ　スカッとジャパン芸能人のスカッと実話をドラマ化１５連発！２時間ＳＰ[字][デ] 20210712",
		"痛快ＴＶ",
		1,
		-1,
		"スカッとジャパン芸能人のスカッと実話をドラマ化１５連発！２時間ＳＰ",
		"2021-07-12",
	)

	do(
		"音楽の日♪中居正広＆安住紳一郎★約８時間の生放送でお届け♪[字][デ] 20210717",
		"音楽の日",
		1,
		-1,
		"中居正広＆安住紳一郎　約８時間の生放送でお届け",
		"2021-07-17",
	)

	do(
		"芸能人格付けチェック　食と芸術の秋３時間スペシャル　浜田爆笑！大物消えた…[デ][字] 20211005",
		"芸能人格付けチェック",
		1,
		-1,
		"食と芸術の秋３時間スペシャル　浜田爆笑！大物消えた…",
		"2021-10-05",
	)

	do(
		"＃ワンチャン賭けとくぅ？　傑作選[字] 20210613",
		"＃ワンチャン賭けとくぅ？　傑作選",
		1,
		-1,
		"",
		"2021-06-13",
	)

	do(
		"♪オトラクションＳＰ♪新企画「厨房の音で料理を当てろ」卓球張本＆SixTONES 20210907",
		"♪オトラクションＳＰ♪新企画「厨房の音で料理を当てろ」卓球張本＆SixTONES",
		1,
		-1,
		"",
		"2021-09-07",
	)
}
