// ./main sendmedia "!429000" test.png                                                                                   0|20:04:13
// ./main msg "!429000" "🐈"

package main

import (
	"github.com/disintegration/imaging"
	"log"
)

func main() {

	// Open the test image.
	src, err := imaging.Open("test.png")
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}
	b := src.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y

	if imgWidth > imgHeight {
		src = imaging.Resize(src, 512, 0, imaging.Lanczos)
	} else {
		src = imaging.Resize(src, 0, 512, imaging.Lanczos)
	}

	// Save the resulting image using JPEG format.
	err = imaging.Save(src, "test2.png")
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}
}
