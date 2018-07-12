package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"path/filepath"
	"time"
)

var supportedFormats = []string{".mp3", ".wav", ".flac"}

func playSong(input Song) (int, error) {
	f, err := os.Open(input.path)
	if err != nil {
		return 0, err
	}
	var s beep.StreamSeekCloser
	var format beep.Format

	switch fileExt := filepath.Ext(input.path); fileExt {
	case ".mp3":
		s, format, err = mp3.Decode(f)
	case ".wav":
		s, format, err = wav.Decode(f)
	case ".flac":
		s, format, err = flac.Decode(f)
	}

	if err != nil {
		return 0, err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
	return int(float32(s.Len()) / float32(format.SampleRate)), nil
}
