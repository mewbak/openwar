/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package game

import (
	"image"
	"time"

	"github.com/andreas-jonsson/openwar/resource"
)

const (
	scrollLeft  = 0x1
	scrollRight = 0x2
	scrollUp    = 0x4
	scrollDown  = 0x8
)

const scrollSpeed = 0.1

type playState struct {
	g *Game
	p *player

	scrollDirection  int
	cameraX, cameraY float64

	ter terrain
	res resource.Resources
}

func NewPlayState(g *Game) GameState {
	ter, _ := newTerrain(g, "HUMAN01")

	return &playState{
		g:   g,
		p:   newPlay(g, humanRace, ter.terrainPalette()),
		res: g.resources,
		ter: ter,
	}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from GameState, args ...interface{}) error {
	s.g.musicPlayer.random(10 * time.Second)

	snd, _ := s.g.soundPlayer.Sound("OREADY.VOC")
	snd.Play(-1, 0, 0)

	return nil
}

func (s *playState) Exit(to GameState) error {
	return nil
}

func (s *playState) Update() error {
	g := s.g
	g.PollAll()

	dt := g.dt
	pos := g.cursorPos
	max := g.renderer.BackBuffer().Bounds().Max
	s.scrollDirection = 0

	if pos.X == 0 {
		s.scrollDirection |= scrollLeft
		s.cameraX -= dt * scrollSpeed
	} else if pos.X == max.X-1 {
		s.scrollDirection |= scrollRight
		s.cameraX += dt * scrollSpeed
	}

	if pos.Y == 0 {
		s.scrollDirection |= scrollUp
		s.cameraY -= dt * scrollSpeed
	} else if pos.Y == max.Y-1 {
		s.scrollDirection |= scrollDown
		s.cameraY += dt * scrollSpeed
	}

	switch {
	case s.scrollDirection == scrollUp|scrollRight:
		g.currentCursor = cursorScrollTopRight
	case s.scrollDirection == scrollDown|scrollRight:
		g.currentCursor = cursorScrollBottomRight
	case s.scrollDirection == scrollDown|scrollLeft:
		g.currentCursor = cursorScrollBottomLeft
	case s.scrollDirection == scrollUp|scrollLeft:
		g.currentCursor = cursorScrollTopLeft
	case s.scrollDirection == scrollUp:
		g.currentCursor = cursorScrollTop
	case s.scrollDirection == scrollRight:
		g.currentCursor = cursorScrollRight
	case s.scrollDirection == scrollDown:
		g.currentCursor = cursorScrollBottom
	case s.scrollDirection == scrollLeft:
		g.currentCursor = cursorScrollLeft
	default:
		g.currentCursor = cursorNormal
	}

	return nil
}

func (s *playState) Render() error {
	mapSize := s.ter.size()
	vp := s.p.hud.viewport()
	cameraPos := image.Point{int(s.cameraX), int(s.cameraY)}
	cameraMax := image.Point{mapSize*16 - (vp.Max.X - vp.Min.X), mapSize*16 - (vp.Max.Y - vp.Min.Y)}

	if cameraPos.X < 0 {
		cameraPos.X = 0
		s.cameraX = 0
	} else if cameraPos.X > cameraMax.X {
		cameraPos.X = cameraMax.X
		s.cameraX = float64(cameraPos.X)
	}

	if cameraPos.Y < 0 {
		cameraPos.Y = 0
		s.cameraY = 0
	} else if cameraPos.Y > cameraMax.Y {
		cameraPos.Y = cameraMax.Y
		s.cameraY = float64(cameraPos.Y)
	}

	s.ter.render(vp, cameraPos)
	s.p.render(s.ter.miniMapImage(), cameraPos)
	return nil
}
