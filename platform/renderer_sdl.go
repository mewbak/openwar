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

package platform

import (
	"image"
	"image/color"
	"image/draw"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type sdlRenderer struct {
	window           *sdl.Window
	backBuffer       *image.RGBA
	internalHWBuffer *sdl.Texture
	internalRenderer *sdl.Renderer
}

func NewRenderer(w, h int, data ...interface{}) (Renderer, error) {
	var (
		r   sdlRenderer
		err error
	)

	r.window, err = sdl.CreateWindow(data[0].(string), sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	r.backBuffer = image.NewRGBA(image.Rect(0, 0, 320, 200))

	renderer, err := sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	r.internalRenderer = renderer

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	renderer.SetLogicalSize(320, 240)
	renderer.SetDrawColor(0, 0, 0, 255)

	r.internalHWBuffer, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, 320, 200)
	if err != nil {
		return nil, err
	}

	sdl.ShowCursor(0)
	return &r, nil
}

func (p *sdlRenderer) ToggleFullscreen() {
	isFullscreen := (p.window.GetFlags() & sdl.WINDOW_FULLSCREEN) != 0
	if isFullscreen {
		p.window.SetFullscreen(0)
	} else {
		//p.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
		p.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
}

func (p *sdlRenderer) Clear() {
	p.internalRenderer.Clear()
	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
}

func (p *sdlRenderer) Present() {
	p.internalHWBuffer.Update(nil, unsafe.Pointer(&p.backBuffer.Pix[0]), p.backBuffer.Stride)
	p.internalRenderer.Copy(p.internalHWBuffer, nil, nil)
	p.internalRenderer.Present()
}

func (p *sdlRenderer) Shutdown() {
	p.window.Destroy()
	p.internalHWBuffer.Destroy()
	p.internalRenderer.Destroy()
}

func (p *sdlRenderer) SetWindowTitle(title string) {
	p.window.SetTitle(title)
}

func (p *sdlRenderer) BackBuffer() draw.Image {
	return p.backBuffer
}

func (p *sdlRenderer) BlitPaletted(dp image.Point, src *image.Paletted) {
	p.BlitImage(dp, src, src.Palette)
}

func (p *sdlRenderer) BlitImage(dp image.Point, src *image.Paletted, pal color.Palette) {
	p.Blit(dp, src, src.Bounds(), pal)
}

func (p *sdlRenderer) Blit(dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	min := sr.Min
	max := sr.Max

	for y, dy := min.Y, 0; y < max.Y; y++ {
		for x, dx := min.X, 0; x < max.X; x++ {
			c := pal[src.ColorIndexAt(x, y)]
			if _, _, _, a := c.RGBA(); a > 0 {
				p.backBuffer.Set(dx+dp.X, dy+dp.Y, c)
			}
			dx++
		}
		dy++
	}
}