package main

import (
	ui "github.com/gizak/termui"
	"fmt"
)

type selectCallback func (int) (int, error);

var songSelectCallback selectCallback

var infoWidget * ui.List
var playlistWidget * ui.List
var scrollerWidget * ui.Gauge
var visualizerWidget * ui.Par
var controlsWidget * ui.Par
var fullSongList []Song

var interfaceSongList []string
var currentSongInterface int = -1
var songPosition = 0
var currentSongLength = 0

func alignInterface() {
	termHeight := ui.TermHeight()
	playlistWidget.Height = termHeight - controlsWidget.Height
	visualizerWidget.Height = termHeight - infoWidget.Height  - controlsWidget.Height - scrollerWidget.Height
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
}
func styleInterface() {
	infoWidget.BorderLabel = "Song info"
	infoWidget.BorderFg = ui.ColorGreen
	playlistWidget.BorderLabel = "PlayList"
	playlistWidget.BorderFg = ui.ColorGreen
	visualizerWidget.BorderLabel = "Visualizer"
	visualizerWidget.BorderFg = ui.ColorGreen
	controlsWidget.BorderFg = ui.ColorGreen
	infoWidget.Height = 6
	controlsWidget.Height = 1
	scrollerWidget.Height = 3
}

func startInterface() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	infoWidget = ui.NewList()
	playlistWidget = ui.NewList()
	scrollerWidget = ui.NewGauge()
	visualizerWidget = ui.NewPar("")
	controlsWidget = ui.NewPar("")
	styleInterface()

	playlistWidget.Items = interfaceSongList;
	
	ui.Body.AddRows(
			ui.NewRow(
				ui.NewCol(6, 0, infoWidget, scrollerWidget, visualizerWidget),
				ui.NewCol(6, 0, playlistWidget)),
			ui.NewRow(
				ui.NewCol(12, 0, controlsWidget)))

	alignInterface()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("timer/1s", func(ui.Event) {
		songPosition++

		if (currentSongLength != 0) {
			scrollerWidget.Percent = int(float32(songPosition) / float32(currentSongLength) * 100)
			ui.Clear()
			ui.Render(ui.Body)
		}
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		songPosition = 0
		var err error
		currentSongLength, err = songSelectCallback(currentSongInterface)
		if err == nil {
			renderCurrentSongInterface(currentSongInterface)
			ui.Clear()
			ui.Render(ui.Body)
		}
	})
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		songUp()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		songDown()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		alignInterface()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Loop()
}

func renderCurrentSongInterface(num int) {
	infoWidget.Items = []string{
		"Artist: " + fullSongList[num].Artist(),
		"Title:  " + fullSongList[num].Title(), 
		"Album:  " + fullSongList[num].Album(),
	}
}

func addSongsInterface(prefix int, inputList []Song) {
	interfaceSongList = make([]string, len(inputList))
	fullSongList = inputList
	for i, v := range inputList {
		if (v.Metadata != nil) {
			interfaceSongList[i] = fmt.Sprintf("[%d] %s - %s", i, v.Artist(), v.Title())
		} else {
			interfaceSongList[i] = fmt.Sprintf("[%d] %s", i, v.path[prefix : ])
		}
	}
	chooseSongInterface(0)
	currentSongInterface = 0
}

func songDown() {
	if (currentSongInterface < len(interfaceSongList) - 1) {
		chooseSongInterface(currentSongInterface + 1)
	}
}

func songUp() {
	if (currentSongInterface > 0) {
		chooseSongInterface(currentSongInterface - 1)
	}
}

func chooseSongInterface(num int) {
	if (currentSongInterface != -1) {
		interfaceSongList[currentSongInterface] = 
		interfaceSongList[currentSongInterface][1: len(interfaceSongList[currentSongInterface]) - 20]
	}
	currentSongInterface = num
	interfaceSongList[num] = fmt.Sprintf("[%s](fg-black,bg-green)", interfaceSongList[num])
}