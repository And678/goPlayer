# GoPlayer
A terminal based audio player

## Features
![screenshot](../assets/screenshot.png)
* Supports mp3, flac, wav formats
* Displays metadata
* Volume controls
* Ability to rewind and fast-forward audio

## Installation

    go get github.com/And678/goPlayer

This will install goPlayer to $GOPATH/bin folder.
Also you can download binaries from 'Releases' tab.

## Usage

To open all audio files in folder: 

    goPlayer /path/to/folder/

To open one specific file: 

    goPlayer /path/to/file.mp3
    
If used without path parameter, goPlayer will assume default music folder: `~/Music/`

## Used libraries

* [termui](https://github.com/gizak/termui/)
* [beep](https://github.com/faiface/beep)
* [tag](https://github.com/dhowden/tag/)

## License
MIT
