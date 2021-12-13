package main

import (
	"image"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: allow to configure this via build tags
const useEmbeddedFiles = true

func OpenFile(path string) (fs.File, error) {
	if useEmbeddedFiles {
		return data.Open(path)
	}

	return os.Open(path)
}

func NewImageFromFile(path string) (*ebiten.Image, image.Image, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	return ebitenutil.NewImageFromReader(file)
}
