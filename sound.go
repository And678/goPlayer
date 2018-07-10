package main
import (
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
	"log"
)

var supportedFormats = []string{".mp3", ".wav", ".flac"}

func playSong(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	s, format, err2 := wav.Decode(f)
	if err2 != nil {
		log.Fatal(err2)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
}