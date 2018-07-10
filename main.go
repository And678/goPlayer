package main

import (
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

func main() {
	var songDir string
	var err error
	if len(os.Args) > 1 {
		songDir = os.Args[1]
	} else {
		songDir, err = homedir.Expand("~/Music/")
		if (err != nil) {
			log.Print("Can't open ~/Music directory")
			os.Exit(1)
		}
	}

	list, err := getSongList(songDir)
	if err != nil {
		log.Print("Can't get song list")
		os.Exit(1)
	}
	addSongsInterface(len(songDir), list)
	songSelectCallback = func (num int) {
		playSong(list[num])
	}
	startInterface()
}