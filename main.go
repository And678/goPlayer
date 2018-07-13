package main

import (
	"github.com/dhowden/tag"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

type Song struct {
	tag.Metadata
	path string
}

func main() {
	var songDir string
	var err error
	if len(os.Args) > 1 {
		songDir = os.Args[1]
	} else {
		songDir, err = homedir.Expand("~/Music/")
		if err != nil {
			log.Fatal("Can't open ~/Music directory")
		}
	}

	fileList, err := getSongList(songDir)
	if err != nil {
		log.Fatal("Can't get song list")
	}
	songs := make([]Song, 0, len(fileList))

	for _, fileName := range fileList {
		currentFile, err := os.Open(fileName)
		if err == nil {
			metadata, _ := tag.ReadFrom(currentFile)
			songs = append(songs, Song{
				Metadata: metadata,
				path:     fileName,
			})
		}
		currentFile.Close()
	}
	if (len(songs) == 0) {
		log.Fatal("Could find any songs to play")
	}
	userInterface, err := NewUi(songs, len(songDir));
	if err != nil {
		log.Fatal(err)
	}
	userInterface.OnSelect = playSong
	userInterface.Start()
	userInterface.Close()
}
