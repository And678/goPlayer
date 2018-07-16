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

type Ui struct {
	infoList		*termui.List
	playList		*termui.List
	scrollerGauge	*termui.Gauge
	visualizer		*termui.Par
	controlsPar		*termui.Par

	songs 			[]Song

	songNum			int

	songSel			int
	songPos			int
	songLen			int

	OnSelect		selectCallback
	OnPause			pauseCallback

	state			uiState
}

func NewUi(songList []Song, pathPrefix int) (*Ui, error) {
	err := termui.Init()
	if err != nil {
		return nil, err
	}
	
	ui := new(Ui)
	ui.songs = songList

	ui.infoList = termui.NewList()
	ui.infoList.BorderLabel = "Song info"
	ui.infoList.BorderFg = termui.ColorGreen
	ui.infoList.Height = 6
	
	ui.playList = termui.NewList()
	ui.playList.BorderLabel = "Playlist"
	ui.playList.BorderFg = termui.ColorGreen
	
	ui.scrollerGauge = termui.NewGauge()
	ui.scrollerGauge.Height = 3
	
	ui.visualizer = termui.NewPar("")
	ui.visualizer.BorderLabel = "Visualizer"
	ui.visualizer.BorderFg = termui.ColorGreen

	ui.controlsPar = termui.NewPar("[p](fg-black,bg-white)[ Pause](fg-black,bg-green) [Esc](fg-black,bg-white)[ Stop](fg-black,bg-green) [left](fg-black,bg-white)[ Forward](fg-black,bg-green) [q](fg-black,bg-white)[ Exit](fg-black,bg-green)")
	ui.controlsPar.BorderFg = termui.ColorGreen
	ui.controlsPar.Border = false
	ui.controlsPar.Height = 1

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, ui.infoList, ui.scrollerGauge, ui.visualizer),
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
				if ui.scrollerGauge.Percent >= 100 {
					ui.songNum++;
					if (ui.songNum >= len(ui.songs)) {
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
	termui.Handle("/sys/kbd/r", func(termui.Event) {
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

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		ui.songDown()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		ui.realign()
	})

	songInter := make([]string, len(ui.songs))
	for i, v := range ui.songs {
		if v.Metadata != nil {
			songInter[i] = fmt.Sprintf("[%d] %s - %s", i + 1, v.Artist(), v.Title())
		} else {
			songInter[i] = fmt.Sprintf("[%d] %s", i + 1, v.path[pathPrefix:])
		}
	}
	ui.playList.Items = songInter
	ui.setSong(0, false)

	return ui, nil
}

func (ui * Ui) Start() {
	termui.Loop()
}

func (ui * Ui) Close() {
	termui.Close()
}


func (ui * Ui) playSong(number int) {
	ui.songPos = 0
	var err error
	ui.songLen, err = ui.OnSelect(ui.songs[number])
	if err == nil {
		ui.state = Playing 
		ui.renderSong()
	}
}

// Rendering

func (ui * Ui) realign() {
	termHeight := termui.TermHeight()
	ui.playList.Height = termHeight - ui.controlsPar.Height
	ui.visualizer.Height = termHeight - ui.infoList.Height - ui.controlsPar.Height - ui.scrollerGauge.Height
	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Clear()
	termui.Render(termui.Body)
}
func (ui * Ui) renderSong() {
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

func (ui * Ui) songDown() {
	if ui.songSel < len(ui.playList.Items) - 1 {
		ui.setSong(ui.songSel + 1, true)
	}
}

func (ui * Ui) songUp() {
	if ui.songSel > 0 {
		ui.setSong(ui.songSel - 1, true)
	}
}

func (ui * Ui) setSong(num int, unset bool) {
	if unset {
		ui.playList.Items[ui.songSel] =
			ui.playList.Items[ui.songSel][1 : len(ui.playList.Items[ui.songSel]) - 20]
	}
	ui.songSel = num
	ui.playList.Items[num] = fmt.Sprintf("[%s](fg-black,bg-green)", ui.playList.Items[num])
}
