package main

import ui "github.com/gizak/termui"

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	infoWidget := ui.NewPar("song/artist window")
	infoWidget.BorderLabel = "Song info"
	infoWidget.BorderFg = ui.ColorCyan

	playlistWidget := ui.NewPar("playList")
	playlistWidget.BorderLabel = "PlayList"
	playlistWidget.BorderFg = ui.ColorYellow

	scrollerWidget := ui.NewPar("")


	visualizerWidget := ui.NewPar("visualizer")
	visualizerWidget.BorderLabel = "Visualizer"
	visualizerWidget.BorderFg = ui.ColorRed

	controlsWidget := ui.NewPar("controlls")
	controlsWidget.BorderFg = ui.ColorGreen

	infoWidget.Height = 6
	controlsWidget.Height = 2
	scrollerWidget.Height = 2
	ui.Body.AddRows(
			ui.NewRow(
				ui.NewCol(6, 0, infoWidget, scrollerWidget, visualizerWidget),
				ui.NewCol(6, 0, playlistWidget)),
			ui.NewRow(
				ui.NewCol(12, 0, controlsWidget)))

	// calculate layout
	ui.Body.Align()
	ui.Render(ui.Body)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		termHeight := ui.TermHeight()
		playlistWidget.Height = termHeight - 2
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})


	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		termHeight := ui.TermHeight()
		playlistWidget.Height = termHeight - 2
		visualizerWidget.Height = termHeight - 4 - 6
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})
	ui.Loop()
}