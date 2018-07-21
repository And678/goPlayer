package main

import (
	"fmt"
	"github.com/gizak/termui"
)

type uiState int

const (
	Stopped uiState = iota
	Playing
	Paused
)

type selectCallback func(Song) (int, error)
type pauseCallback func(bool)
type seekCallback func(int) error
type volumeCallback func(int)

type Ui struct {
	infoList      *termui.List
	playList      *termui.List
	scrollerGauge *termui.Gauge
	volumeGauge   *termui.Gauge
	controlsPar   *termui.Par

	songs     []Song
	songNames []string

	volume int

	songNum int

	songSel int
	songPos int
	songLen int

	OnSelect selectCallback
	OnPause  pauseCallback
	OnSeek   seekCallback
	OnVolume volumeCallback

	state uiState
}

func NewUi(songList []Song, pathPrefix int) (*Ui, error) {
	err := termui.Init()
	if err != nil {
		return nil, err
	}

	ui := new(Ui)

	ui.volume = 100

	ui.songs = songList

	ui.infoList = termui.NewList()
	ui.infoList.BorderLabel = "Song info"
	ui.infoList.BorderFg = termui.ColorGreen

	ui.playList = termui.NewList()
	ui.playList.BorderLabel = "Playlist"
	ui.playList.BorderFg = termui.ColorGreen

	ui.scrollerGauge = termui.NewGauge()
	ui.scrollerGauge.BorderLabel = "Stopped"
	ui.scrollerGauge.Height = 3

	ui.volumeGauge = termui.NewGauge()
	ui.volumeGauge.BorderLabel = "Volume"
	ui.volumeGauge.Height = 3
	ui.volumeGauge.Percent = ui.volume

	ui.controlsPar = termui.NewPar(
		"[ Enter ](fg-black,bg-white)[ Select ](fg-black,bg-green) " +
			"[ p ](fg-black,bg-white)[ Pause ](fg-black,bg-green) " +
			"[Esc](fg-black,bg-white)[ Stop ](fg-black,bg-green) " +
			"[Right](fg-black,bg-white)[ +10s ](fg-black,bg-green) " +
			"[Left](fg-black,bg-white)[ -10s ](fg-black,bg-green) " +
			"[ q ](fg-black,bg-white)[ Exit ](fg-black,bg-green) ")
	ui.controlsPar.Border = false
	ui.controlsPar.Height = 1

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, ui.infoList, ui.scrollerGauge, ui.volumeGauge),
			termui.NewCol(6, 0, ui.playList)),
		termui.NewRow(
			termui.NewCol(12, 0, ui.controlsPar)))

	ui.realign()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/p", func(termui.Event) {
		if ui.state == Playing {
			ui.OnPause(true)
			ui.state = Paused
		} else {
			ui.OnPause(false)
			ui.state = Playing

		}
		ui.renderSong()
	})
	termui.Handle("timer/1s", func(termui.Event) {
		if ui.state == Playing {
			ui.songPos++
			if ui.songLen != 0 {
				ui.scrollerGauge.Percent = int(float32(ui.songPos) / float32(ui.songLen) * 100)
				ui.scrollerGauge.Label = fmt.Sprintf("%d:%.2d / %d:%.2d", ui.songPos/60, ui.songPos%60, ui.songLen/60, ui.songLen%60)
				if ui.scrollerGauge.Percent >= 100 {
					ui.songNum++
					if ui.songNum >= len(ui.songs) {
						ui.songNum = 0
					}
					ui.playSong(ui.songNum)
				}
				termui.Clear()
				termui.Render(termui.Body)
			}
		} else if ui.state == Stopped {
			ui.songPos = 0
		}
	})

	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		ui.songPos += 10
		ui.OnSeek(ui.songPos)
	})

	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		ui.songPos -= 10
		if ui.songPos < 0 {
			ui.songPos = 0
		}
		ui.OnSeek(ui.songPos)
	})

	termui.Handle("/sys/kbd/<escape>", func(termui.Event) {
		ui.playSong(ui.songNum)
		ui.OnPause(true)
		ui.state = Stopped
		ui.scrollerGauge.Percent = 0
		ui.renderSong()
	})

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		ui.songNum = ui.songSel
		ui.playSong(ui.songNum)
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		ui.songUp()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/=", func(termui.Event) {
		ui.volumeUp()
	})

	termui.Handle("/sys/kbd/+", func(termui.Event) {
		ui.volumeUp()
	})

	termui.Handle("/sys/kbd/-", func(termui.Event) {
		ui.volumeDown()
	})

	termui.Handle("/sys/kbd/_", func(termui.Event) {
		ui.volumeDown()
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		ui.songDown()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		ui.realign()
	})

	ui.songNames = make([]string, len(ui.songs))
	for i, v := range ui.songs {
		if v.Metadata != nil {
			ui.songNames[i] = fmt.Sprintf("[%d] %s - %s", i+1, v.Artist(), v.Title())
		} else {
			ui.songNames[i] = fmt.Sprintf("[%d] %s", i+1, v.path[pathPrefix:])
		}
	}
	ui.playList.Items = ui.songNames
	ui.setSong(0, false)

	return ui, nil
}

func (ui *Ui) Start() {
	termui.Loop()
}

func (ui *Ui) Close() {
	termui.Close()
}

func (ui *Ui) playSong(number int) {
	ui.songPos = 0
	var err error
	ui.songLen, err = ui.OnSelect(ui.songs[number])
	if err == nil {
		ui.state = Playing
		ui.renderSong()
	}
}

// Rendering

func (ui *Ui) realign() {
	termHeight := termui.TermHeight()
	ui.playList.Height = termHeight - ui.controlsPar.Height
	ui.infoList.Height = termHeight - ui.controlsPar.Height - ui.scrollerGauge.Height - ui.volumeGauge.Height
	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Clear()
	termui.Render(termui.Body)
}

func (ui *Ui) renderSong() {
	var status string
	switch ui.state {
	case Playing:
		status = "[(Playing)](fg-green)"
	case Paused:
		status = "[(Paused)](fg-yellow)"
	case Stopped:
		status = "[(Stopped)](fg-red)"
	}

	ui.infoList.Items = []string{
		"Artist: " + ui.songs[ui.songSel].Artist(),
		"Title:  " + ui.songs[ui.songSel].Title(),
		"Album:  " + ui.songs[ui.songSel].Album(),
		status,
	}
	termui.Clear()
	termui.Render(termui.Body)
}

//Song selection

func (ui *Ui) songDown() {
	if ui.songSel < len(ui.songNames)-1 {
		ui.setSong(ui.songSel+1, true)
	}
}

func (ui *Ui) songUp() {
	if ui.songSel > 0 {
		ui.setSong(ui.songSel-1, true)
	}
}

func (ui *Ui) volumeUp() {
	if ui.volume < 100 {
		ui.volume += 5
	}
	ui.volumeGauge.Percent = ui.volume
	ui.OnVolume(ui.volume)
	termui.Clear()
	termui.Render(termui.Body)
}

func (ui *Ui) volumeDown() {
	if ui.volume > 0 {
		ui.volume -= 5
	}
	ui.volumeGauge.Percent = ui.volume
	ui.OnVolume(ui.volume)
	termui.Clear()
	termui.Render(termui.Body)
}

func (ui *Ui) setSong(num int, unset bool) {
	skip := 0
	for num-skip >= ui.playList.Height-2 {
		skip += ui.playList.Height - 2
	}
	if unset {
		ui.songNames[ui.songSel] = ui.songNames[ui.songSel][1 : len(ui.songNames[ui.songSel])-20]
	}
	ui.songSel = num
	ui.songNames[num] = fmt.Sprintf("[%s](fg-black,bg-green)", ui.songNames[num])
	ui.playList.Items = ui.songNames[skip:]
}
