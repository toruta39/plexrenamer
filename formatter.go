package main

import (
	"fmt"
)

func PlexFormat(programInfo ProgramInfo, ext string) string {
	var descriptor string

	if programInfo.Episode == -1 {
		descriptor = programInfo.Date.Format("2006-01-02")
	} else {
		descriptor = fmt.Sprintf("S%02dE%02d", programInfo.Season, programInfo.Episode)
	}

	if programInfo.Subtitle != "" {
		descriptor += " " + programInfo.Subtitle
	}

	return fmt.Sprintf("%s/Season %02d/%s %s.%s",
		programInfo.Title,
		programInfo.Season,
		programInfo.Title,
		descriptor,
		ext,
	)
}
