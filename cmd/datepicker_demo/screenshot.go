package main

import (
	"image/png"
	"log"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"github.com/sdassow/atomic"
)

func makeScreenshot(w fyne.Window) {
	dir := "."
	basename := "screenshot"
	filename := basename + ".png"

	img := w.Canvas().Capture()
	f, err := os.CreateTemp(dir, basename+"-*.png")
	if err != nil {
		log.Printf("failed to create temporary file: %v", err)
		return
	}

	if err := png.Encode(f, img); err != nil {
		log.Printf("failed to encode image: %v", err)
		return
	}

	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v", err)
		return
	}

	if err := atomic.ReplaceFile(f.Name(), path.Join(dir, filename)); err != nil {
		log.Printf("failed to rename file: %v", err)
		return
	}

	log.Printf("screenshot saved: %s", path.Join(dir, filename))
}
