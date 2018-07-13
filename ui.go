package main

import (
	"fmt"
	"github.com/gizak/termui"
)

type selectCallback func(Song) (int, error)

type Ui struct {
	infoList		*termui.List
	playList		*termui.List
	scrollerGauge	*termui.Gauge
	visualizer		*termui.Par
	controlsPar		*termui.Par

	songs 			[]Song
	songNum	int
	songPos	int
	songLen	int

	OnSelect		selectCallback
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
	ui.playList.BorderLabel = "PlayList"
	ui.playList.BorderFg = termui.ColorGreen
	
	ui.scrollerGauge = termui.NewGauge()
	ui.scrollerGauge.Height = 3
	
	ui.visualizer = termui.NewPar("")
	ui.visualizer.BorderLabel = "Visualizer"
	ui.visualizer.BorderFg = termui.ColorGreen

	ui.controlsPar = termui.NewPar("[p](fg-black,bg-white)[ Pause](fg-black,bg-green) [Esc](fg-black,bg-white)[ Stop](fg-black,bg-green) [q](fg-black,bg-white)[ Exit](fg-black,bg-green) [q](fg-black,bg-white)[ Exit](fg-black,bg-green)")
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
	termui.Handle("timer/1s", func(termui.Event) {
		ui.songPos++
		if ui.songLen != 0 {
			ui.scrollerGauge.Percent = int(float32(ui.songPos) / float32(ui.songLen) * 100)
			termui.Clear()
			termui.Render(termui.Body)
		}
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		ui.songPos = 0
		var err error
		ui.songLen, err = ui.OnSelect(ui.songs[ui.songNum])
		if err == nil {
			ui.renderSong(ui.songNum)
			termui.Clear()
			termui.Render(termui.Body)
		}
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
			songInter[i] = fmt.Sprintf("[%d] %s - %s", i, v.Artist(), v.Title())
		} else {
			songInter[i] = fmt.Sprintf("[%d] %s", i, v.path[pathPrefix:])
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

func (ui * Ui) realign() {
	termHeight := termui.TermHeight()
	ui.playList.Height = termHeight - ui.controlsPar.Height
	ui.visualizer.Height = termHeight - ui.infoList.Height - ui.controlsPar.Height - ui.scrollerGauge.Height
	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Clear()
	termui.Render(termui.Body)
}
func (ui * Ui) renderSong(num int) {
	ui.infoList.Items = []string{
		"Artist: " + ui.songs[num].Artist(),
		"Title:  " + ui.songs[num].Title(),
		"Album:  " + ui.songs[num].Album(),
	}
}

func (ui * Ui) songDown() {
	if ui.songNum < len(ui.playList.Items) - 1 {
		ui.setSong(ui.songNum + 1, true)
	}
}

func (ui * Ui) songUp() {
	if ui.songNum > 0 {
		ui.setSong(ui.songNum - 1, true)
	}
}

func (ui * Ui) setSong(num int, unset bool) {
	if unset {
		ui.playList.Items[ui.songNum] =
			ui.playList.Items[ui.songNum][1 : len(ui.playList.Items[ui.songNum]) - 20]
	}
	ui.songNum = num
	ui.playList.Items[num] = fmt.Sprintf("[%s](fg-black,bg-green)", ui.playList.Items[num])
}
