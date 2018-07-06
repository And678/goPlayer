package main

import (
	ui "github.com/gizak/termui"
)

var infoWidget * ui.Par
var playlistWidget * ui.Par
var scrollerWidget * ui.Par
var visualizerWidget * ui.Par
var controlsWidget * ui.Par

func alignInterface() {
	termHeight := ui.TermHeight()
	playlistWidget.Height = termHeight - controlsWidget.Height
	visualizerWidget.Height = termHeight - infoWidget.Height  - controlsWidget.Height - scrollerWidget.Height
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
}
func styleInterface() {
	infoWidget.BorderLabel = "Song info"
	infoWidget.BorderFg = ui.ColorCyan
	playlistWidget.BorderLabel = "PlayList"
	playlistWidget.BorderFg = ui.ColorYellow
	visualizerWidget.BorderLabel = "Visualizer"
	visualizerWidget.BorderFg = ui.ColorRed
	controlsWidget.BorderFg = ui.ColorGreen
	infoWidget.Height = 6
	controlsWidget.Height = 1
	scrollerWidget.Height = 2
}

func startInterface() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	infoWidget = ui.NewPar("")
	playlistWidget = ui.NewPar("")
	scrollerWidget = ui.NewPar("")
	visualizerWidget = ui.NewPar("")
	controlsWidget = ui.NewPar("")
	styleInterface()
	
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

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		alignInterface()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Loop()
}