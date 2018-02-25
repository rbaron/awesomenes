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
    SCREEN_WIDTH * 3,
    SCREEN_HEIGHT * 3,
    sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}

  renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

  texture, err := renderer.CreateTexture(
    sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, SCREEN_WIDTH, SCREEN_HEIGHT)

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

func (tv *TV) SetFrame(pixels []byte) {
  //log.Printf("WIll set frame")
  tv.texture.Update(nil, pixels, SCREEN_WIDTH)
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
