package main

import (
	"github.com/mitchellh/go-homedir"
	"github.com/dhowden/tag"
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
		if (err != nil) {
			log.Print("Can't open ~/Music directory")
			os.Exit(1)
		}
	}

	fileList, err := getSongList(songDir)
	if err != nil {
		log.Print("Can't get song list")
		os.Exit(1)
	}
	songs := make([]Song, 0, len(fileList))

	for _, fileName := range fileList {
		currentFile, _ := os.Open(fileName) // TODO: handle the error
		metadata, _ := tag.ReadFrom(currentFile)
		songs = append(songs, Song{
			Metadata: metadata,
			path: fileName,
		})
		currentFile.Close()
	}

	addSongsInterface(len(songs), songs)
	songSelectCallback = func (num int) {
		playSong(fileList[num])
	}
	startInterface()
}