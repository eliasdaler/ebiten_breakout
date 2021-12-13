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
