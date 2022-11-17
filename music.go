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
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Music struct {
	f      fs.File
	loop   *audio.InfiniteLoop
	player *audio.Player
}

func NewMusicFromFile(audioContext *audio.Context, path string) (*Music, error) {
	m := &Music{}

	f, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	m.f = f

	s, err := vorbis.Decode(audioContext, m.f)
	if err != nil {
		return nil, err
	}

	m.loop = audio.NewInfiniteLoop(s, s.Length())

	p, err := audio.NewPlayer(audioContext, m.loop)
	if err != nil {
		return nil, err
	}
	m.player = p

	return m, nil
}

func (m *Music) Play() {
	m.player.Play()
}

func (m *Music) Close() error {
	if err := m.player.Close(); err != nil {
		return err
	}
	if err := m.f.Close(); err != nil {
		return err
	}
	return nil
}
