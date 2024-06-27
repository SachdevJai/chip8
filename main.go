package main

import (
	"chip8/cpu"
	"os"

	sdl "github.com/veandco/go-sdl2/sdl"
)

const WIDTH int32 = 64
const HEIGHT int32 = 32
const FPS uint32 = 500

func main() {

	fileName := "roms/" + string(os.Args[1]) + ".ch8"

	var sizeModifier int32 = 20

	var c8 cpu.Chip8 = cpu.Init()

	loadErr := c8.LoadProgram(fileName)

	if loadErr != nil {
		panic(loadErr)
	}

	// Init SDL
	var sdlErr = sdl.Init(sdl.INIT_EVERYTHING)
	if sdlErr != nil {
		panic(sdlErr)
	}
	defer sdl.Quit()

	// Create window
	window, windowErr := sdl.CreateWindow("c", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WIDTH*sizeModifier, HEIGHT*sizeModifier, sdl.WINDOW_SHOWN)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	// Create renderer
	renderer, rendererErr := sdl.CreateRenderer(window, -1, 0)
	if rendererErr != nil {
		panic(rendererErr)
	}
	defer renderer.Destroy()

	shouldrun := true

	// Render Loop
	for shouldrun {

		c8.Cycle()

		if c8.Draw() {
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.Clear()

			displayBuffer := c8.Buffer()

			for i := 0; i < len(displayBuffer); i++ {
				for j := 0; j < len(displayBuffer[i]); j++ {

					// fmt.Printf("%d", displayBuffer[i][j])

					if displayBuffer[i][j] == 1 {
						renderer.SetDrawColor(255, 255, 255, 255)
					} else {
						renderer.SetDrawColor(0, 0, 0, 255)
					}

					renderer.FillRect(
						&sdl.Rect{
							X: int32(j) * sizeModifier,
							Y: int32(i) * sizeModifier,
							W: sizeModifier,
							H: sizeModifier,
						})
				}
			}
			renderer.Present()
		}

		event := sdl.PollEvent()

		switch et := event.(type) {
		case *sdl.QuitEvent:
			shouldrun = false
		case *sdl.KeyboardEvent:
			if et.Type == sdl.KEYUP {
				switch et.Keysym.Sym {
				case sdl.K_1:
					c8.KeyRelease(0x1)
				case sdl.K_2:
					c8.KeyRelease(0x2)
				case sdl.K_3:
					c8.KeyRelease(0x3)
				case sdl.K_4:
					c8.KeyRelease(0xC)
				case sdl.K_q:
					c8.KeyRelease(0x4)
				case sdl.K_w:
					c8.KeyRelease(0x5)
				case sdl.K_e:
					c8.KeyRelease(0x6)
				case sdl.K_r:
					c8.KeyRelease(0xD)
				case sdl.K_a:
					c8.KeyRelease(0x7)
				case sdl.K_s:
					c8.KeyRelease(0x8)
				case sdl.K_d:
					c8.KeyRelease(0x9)
				case sdl.K_f:
					c8.KeyRelease(0xE)
				case sdl.K_z:
					c8.KeyRelease(0xA)
				case sdl.K_x:
					c8.KeyRelease(0x0)
				case sdl.K_c:
					c8.KeyRelease(0xB)
				case sdl.K_v:
					c8.KeyRelease(0xF)
				}
			} else if et.Type == sdl.KEYDOWN {
				switch et.Keysym.Sym {
				case sdl.K_1:
					c8.KeyPress(0x1)
				case sdl.K_2:
					c8.KeyPress(0x2)
				case sdl.K_3:
					c8.KeyPress(0x3)
				case sdl.K_4:
					c8.KeyPress(0xC)
				case sdl.K_q:
					c8.KeyPress(0x4)
				case sdl.K_w:
					c8.KeyPress(0x5)
				case sdl.K_e:
					c8.KeyPress(0x6)
				case sdl.K_r:
					c8.KeyPress(0xD)
				case sdl.K_a:
					c8.KeyPress(0x7)
				case sdl.K_s:
					c8.KeyPress(0x8)
				case sdl.K_d:
					c8.KeyPress(0x9)
				case sdl.K_f:
					c8.KeyPress(0xE)
				case sdl.K_z:
					c8.KeyPress(0xA)
				case sdl.K_x:
					c8.KeyPress(0x0)
				case sdl.K_c:
					c8.KeyPress(0xB)
				case sdl.K_v:
					c8.KeyPress(0xF)
				case sdl.K_ESCAPE:
					shouldrun = false
				}
			}
		}

		// fps
		sdl.Delay(1000 / FPS)
	}
}
