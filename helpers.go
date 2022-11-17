// Copyright (c) 2022 Elias Daler
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
