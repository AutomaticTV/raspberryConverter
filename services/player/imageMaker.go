package player

import (
	"bytes"
	"errors"
	"image"
	"os"

	"github.com/fogleman/gg"
	"github.com/gobuffalo/packr"
)

// w represents the width in pixels of the image
const w = 1920

// h represents the height in pixels of the image
const h = 1080

// box packs static assets (base image and font) required to create the IP image
var box = packr.NewBox("./assets")

// initImageMaker loads the necessary resources to create the desired images,
// and can be reused in the different executions
func initImageMaker() (*gg.Context, error) {
	// CREATE THE DESTINATION FOLDER IF NEEDED
	imageMaker := gg.NewContext(w, h)
	err := os.MkdirAll(destinationPath, 0777)
	if err != nil {
		return imageMaker, errors.New("Error creating folder " + destinationPath + ": " + err.Error())
	}
	// INIT THE imageMaker
	// LOAD FONT
	fileName := "font.ttf"
	tmpFilepath := destinationPath + fileName
	errorLoadingFont := func(err error) (*gg.Context, error) {
		return imageMaker, errors.New("Error loading font: " + err.Error())
	}
	fontBytes, err := box.Find(fileName)
	if err != nil {
		return errorLoadingFont(err)
	}
	f, err := os.Create(tmpFilepath)
	defer f.Close()
	if err != nil {
		return errorLoadingFont(err)
	}
	if _, err := f.Write(fontBytes); err != nil {
		return errorLoadingFont(err)
	}
	if err := imageMaker.LoadFontFace(tmpFilepath, 70); err != nil {
		return errorLoadingFont(err)
	}
	return imageMaker, nil
}

// makeImage stores a image at destinationFile, this new image is based on services/player/assets/bg.png with a label text placed in the middle of the image
func makeImage(label string) error {
	imageMaker, err := initImageMaker()
	if err != nil {
		return errors.New("Error initializig imageMaker: " + err.Error())
	}
	// FONT COLOR
	imageMaker.SetRGB(1, 1, 1)
	// LOAD BASE IMAGE
	imgBytes, err := box.Find("bg.png")
	if err != nil {
		return errors.New("Error geting background image to generate the display IP image: " + err.Error())
	}
	imgReader := bytes.NewReader(imgBytes)
	baseImage, _, err := image.Decode(imgReader)
	if err != nil {
		return errors.New("Error decoding the image: " + err.Error())
	}
	// SET BASE IMAGE AS BACKGROUND
	imageMaker.DrawImage(baseImage, 0, 0)
	// PRINT THE LABEL IN THE CENTER OF THE BACKGROUND
	imageMaker.DrawStringAnchored(label, w/2, h/2, 0.5, 0.5)
	// SAVE THE NEW IMAGE
	imageMaker.Clip()
	return imageMaker.SavePNG(destinationFile)
}
