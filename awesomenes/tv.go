package awesomenes

import (
	//"log"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

type TV struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func MakeTV() *TV {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(
		"Awesomenes",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		SCREEN_WIDTH*3,
		SCREEN_HEIGHT*3,
		sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, SCREEN_WIDTH, SCREEN_HEIGHT)

	if err != nil {
		panic(err)
	}

	renderer.SetLogicalSize(SCREEN_WIDTH, SCREEN_HEIGHT)

	return &TV{
		window:   window,
		renderer: renderer,
		texture:  texture,
	}
}

func (tv *TV) UpdateInputState(ctrlr *Controller) {
	for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
		switch evt.(type) {
		case *sdl.KeyboardEvent:
			tv.handleKBDEvevent(ctrlr, evt.(*sdl.KeyboardEvent))
		}
	}
}

func (tv *TV) handleKBDEvevent(ctrlr *Controller, evt *sdl.KeyboardEvent) {
	if evt.Repeat != 0 {
		return
	}

	fn := ctrlr.PushButton

	if evt.Type == sdl.KEYUP {
		fn = ctrlr.ReleaseButton
	}

	switch evt.Keysym.Sym {
	case sdl.K_RETURN:
		fn(CONTROLLER_BUTTONS_START)
	case sdl.K_RSHIFT:
		fn(CONTROLLER_BUTTONS_SELECT)
	case sdl.K_a:
		fn(CONTROLLER_BUTTONS_A)
	case sdl.K_s:
		fn(CONTROLLER_BUTTONS_B)
	case sdl.K_UP:
		fn(CONTROLLER_BUTTONS_UP)
	case sdl.K_RIGHT:
		fn(CONTROLLER_BUTTONS_RIGHT)
	case sdl.K_DOWN:
		fn(CONTROLLER_BUTTONS_DOWN)
	case sdl.K_LEFT:
		fn(CONTROLLER_BUTTONS_LEFT)
	}
}

func (tv *TV) SetFrame(pixels []byte) {
	tv.texture.Update(nil, pixels, SCREEN_WIDTH*4)
}

func (tv *TV) ShowPixels() {
	tv.renderer.Clear()
	tv.renderer.Copy(tv.texture, nil, nil)
	tv.renderer.Present()
}

func (tv *TV) Cleanup() {
	tv.window.Destroy()
	sdl.Quit()
}
