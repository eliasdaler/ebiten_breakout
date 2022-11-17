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
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type Sound struct {
	f      fs.File
	stream *wav.Stream
	player *audio.Player
}

func NewSoundFromFile(audioContext *audio.Context, path string) (*Sound, error) {
	s := &Sound{}

	f, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	s.f = f

	ws, err := wav.Decode(audioContext, s.f)
	if err != nil {
		return nil, err
	}
	s.stream = ws

	p, err := audioContext.NewPlayer(ws)
	if err != nil {
		return nil, err
	}
	s.player = p

	return s, nil
}

func (s *Sound) Play() {
	s.player.Rewind()
	s.player.Play()
}

func (s *Sound) Close() error {
	if err := s.player.Close(); err != nil {
		return err
	}
	if err := s.f.Close(); err != nil {
		return err
	}
	return nil
}
