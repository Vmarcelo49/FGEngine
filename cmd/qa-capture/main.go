//go:build fgengineqa

package main

import (
	"fgengine/scene"
	"flag"
	"log"
)

func main() {
	outDir := flag.String("out", "artifacts/gameplay", "directory for captured PNG frames")
	flag.Parse()
	if err := scene.CaptureRosterFrames(*outDir); err != nil {
		log.Fatal(err)
	}
}
