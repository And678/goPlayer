package main

import (
	"bytes"
	"github.com/eliukblau/pixterm/ansimage"
	"image/color"
)

func drawImage(x, y int, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	image, err := ansimage.NewScaledFromReader(reader, x, y, color.Black, ansimage.ScaleModeFit, ansimage.NoDithering)
	if err != nil {
		return "", err
	}

	return image.Render(), nil
}
