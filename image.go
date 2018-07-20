package main

import (
	"github.com/eliukblau/pixterm/ansimage"
	"bytes"
	"image/color"
	)

func drawImage(x,y int, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	image, err := ansimage.NewScaledFromReader(reader, x, y, color.Black, ansimage.ScaleModeFit, ansimage.NoDithering)
	if (err != nil) {
		return "", err
	}

	return image.Render(), nil
}