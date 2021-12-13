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
