package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"path/filepath"
	"time"
)

var supportedFormats = []string{".mp3", ".wav", ".flac"}
var mainCtrl *beep.Ctrl
var s beep.StreamSeekCloser
var format beep.Format
var volume = &effects.Volume{
	Base: 2,
}

func playSong(input Song) (int, error) {
	f, err := os.Open(input.path)
	if err != nil {
		return 0, err
	}

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
	volume.Streamer = s
	mainCtrl = &beep.Ctrl{Streamer: volume}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(mainCtrl)
	return int(float32(s.Len()) / float32(format.SampleRate)), nil
}

func pauseSong(state bool) {
	speaker.Lock()
	mainCtrl.Paused = state
	speaker.Unlock()
}

func seek(pos int) error {
	speaker.Lock()
	err := s.Seek(pos * int(format.SampleRate))
	speaker.Unlock()
	return err
}

func setVolue(percent int) {
	if percent == 0 {
		volume.Silent = true
	} else {
		volume.Silent = false
		volume.Volume = -float64(100-percent) / 100.0 * 5
	}
}
