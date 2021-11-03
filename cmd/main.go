package main

import (
	"flag"
	"fmt"

	plexrenamer "github.com/toruta39/plex-renamer"
)

func main() {
	fromDir := flag.String("from", "", "Source directory")
	toDir := flag.String("to", "", "Target directory")
	dryrun := flag.Bool("dryrun", false, "Dry run")
	flag.Parse()

	if *fromDir == "" || *toDir == "" {
		flag.PrintDefaults()
		return
	}

	toFiles, err := plexrenamer.ScanDir(*fromDir, *toDir, *dryrun)
	if err != nil {
		panic(err)
	}

	for _, file := range toFiles {
		fmt.Printf("Moved: %s\n", file)
	}
}
