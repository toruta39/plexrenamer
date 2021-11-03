package plexrenamer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const delimiter = "|"

type ProgramInfo struct {
	Title    string
	Season   int
	Episode  int
	Subtitle string
	Date     time.Time
}

func dedupeDelimiter(input string) string {
	oldSegments := strings.Split(input, delimiter)
	newSegments := make([]string, 0)
	for _, v := range oldSegments {
		v = strings.TrimSpace(v)
		if v != "" {
			newSegments = append(newSegments, v)
		}
	}
	return strings.Join(newSegments, delimiter)
}

func predelimit(input string) string {
	re := regexp.MustCompile("★|☆")
	output := re.ReplaceAllString(input, delimiter)

	knownTitles := []string{
		"はやドキ！",
	}

	for _, title := range knownTitles {
		output = strings.ReplaceAll(output, title, title+delimiter)
	}

	return output
}

func sanitizeSlots(input string) string {
	re := regexp.MustCompile("＜[^＜＞]+＞")
	return re.ReplaceAllString(input, delimiter)
}

func sanitizeFlags(input string) string {
	re := regexp.MustCompile(`\[[新|再|字|二|多|解|吹|初|終|デ]\]`)
	return re.ReplaceAllString(input, delimiter)
}

func normalizeTitle(input string) string {
	re := regexp.MustCompile(`^(アニメ|ドラマシリーズ|ドラマ２５|ドラマ)　?`)
	input = re.ReplaceAllString(input, "")
	input = strings.Trim(input, "『』")
	return input
}

func normalizeSubtitle(input string) string {
	return strings.ReplaceAll(input, delimiter, "　")
}

func normalizeNumber(input string) string {
	input = strings.ReplaceAll(input, "０", "0")
	input = strings.ReplaceAll(input, "１", "1")
	input = strings.ReplaceAll(input, "２", "2")
	input = strings.ReplaceAll(input, "３", "3")
	input = strings.ReplaceAll(input, "４", "4")
	input = strings.ReplaceAll(input, "５", "5")
	input = strings.ReplaceAll(input, "６", "6")
	input = strings.ReplaceAll(input, "７", "7")
	input = strings.ReplaceAll(input, "８", "8")
	input = strings.ReplaceAll(input, "９", "9")
	return input
}

func matchDate(input string) (time.Time, string, error) {
	re := regexp.MustCompile(`(\d{4})-?(\d{2})-?(\d{2})`)

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return time.Unix(0, 0), input, nil
	}

	d, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", matches[1], matches[2], matches[3]))
	if err != nil {
		return time.Now(), input, err
	}

	output := dedupeDelimiter(strings.ReplaceAll(input, matches[0], delimiter))
	return d, output, nil
}

func matchEpisode(input string) (int, string, error) {
	re := regexp.MustCompile(`[#＃第（]([0-9０-９]{1,2})[）話]?`)

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return -1, input, nil
	}

	episode, err := strconv.Atoi(normalizeNumber(matches[1]))
	if err != nil {
		return -1, input, err
	}

	output := dedupeDelimiter(strings.ReplaceAll(input, matches[0], delimiter))
	return episode, output, nil
}

func matchSubtitle(input string) (string, string, error) {
	re := regexp.MustCompile(`[【「♪](.{2,})[♪」】]$`)

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return "", input, nil
	}

	output := dedupeDelimiter(strings.ReplaceAll(input, matches[0], delimiter))

	return normalizeSubtitle(matches[1]), output, nil
}

func matchTitle(input string) (string, string, string, error) {
	segments := strings.Split(input, delimiter)
	if segments == nil {
		return "", "", input, fmt.Errorf("Title not found")
	}

	segments[0] = normalizeTitle(segments[0])

	// use []rune to correctly count length in unicode chars
	if len([]rune(segments[0])) < 20 {
		return segments[0], normalizeSubtitle(strings.Join(segments[1:], "　")), "", nil
	}

	re := regexp.MustCompile(`\s|　`)

	segments = append(re.Split(segments[0], -1), segments[1:]...)
	fallbackSubtitle := ""
	if len(segments) > 1 {
		fallbackSubtitle = strings.Join(segments[1:], "　")
	}
	return segments[0], normalizeSubtitle(fallbackSubtitle), "", nil
}

func Parse(input string) (*ProgramInfo, error) {
	result := &ProgramInfo{
		Season:  1,
		Episode: -1,
	}

	input = dedupeDelimiter(predelimit(sanitizeSlots(sanitizeFlags(input))))

	date, input, err := matchDate(input)
	if err != nil {
		return nil, err
	}
	result.Date = date

	episode, input, err := matchEpisode(input)
	if err != nil {
		return nil, err
	}
	result.Episode = episode

	subtitle, input, err := matchSubtitle(input)
	if err != nil {
		return nil, err
	}
	result.Subtitle = subtitle

	title, fallbackSubtitle, input, err := matchTitle(input)
	if err != nil {
		return nil, err
	}
	result.Title = title

	if result.Subtitle == "" {
		result.Subtitle = fallbackSubtitle
	}

	return result, nil
}
