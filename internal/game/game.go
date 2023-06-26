package game

import (
	"image"

	"github.com/tahaontech/go_ssr_game_engine/internal/frame"
)

type Command struct {
	Dir int `json:"dir"`
}

type GameObj struct {
	Dir           int // 0 rotate clockwise, 1 unclockwise
	State         int // 1 - 365 rotate degree
	FrameRenderer *frame.FrameControl
}

func NewGame(path string) (*GameObj, error) {
	fc, err := frame.NewFrameControl(path)
	if err != nil {
		return nil, err
	}
	return &GameObj{
		Dir:           0,
		State:         1,
		FrameRenderer: fc,
	}, nil
}

func (g *GameObj) Reset() {
	g.State = 1
	g.Dir = 0
}

func (g *GameObj) Update() {
	dir := g.Dir
	if dir == 0 {
		g.State = g.State + 1
		if g.State == 360 {
			g.State = 0
		}
	} else {
		g.State = g.State - 1
		if g.State == -1 {
			g.State = 360
		}
	}
}

func (g *GameObj) GetFrame() image.Image {
	return g.FrameRenderer.DrawFrame(g.State)
}
