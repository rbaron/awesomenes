package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rbaron/awesomenes/awesomenes"
)

const (
	FRAME_RATE   = 60
	FRAME_CYCLES = awesomenes.CPU_FREQ / FRAME_RATE

	// All timings in nanoseconds
	FRAME_TIME = 1e9 / float64(FRAME_RATE)
)

func main() {
	tv := awesomenes.MakeTV()

	if len(os.Args) != 2 {
		fmt.Println("Usage: awesomenes ROM_PATH")
		os.Exit(2)
	}

	rom := awesomenes.ReadROM(os.Args[1])
  mapper := awesomenes.MakeMapper(rom)

	ppu := awesomenes.MakePPU(nil, rom, mapper)
	ppu.Reset()

	controller := awesomenes.MakeController()
	cpuAddrSpace := awesomenes.MakeCPUAddrSpace(rom, ppu, controller, mapper)
	cpu := awesomenes.MakeCPU(cpuAddrSpace)

	ppu.CPU = cpu
	ppu.TV = tv

	var cpuCycles, newCycles, ppuCycles int

	t0 := time.Now().UnixNano()

	for {

		newCycles = cpu.Run()
		cpuCycles += newCycles

		for ppuCycles = 0; ppuCycles < 3*newCycles; ppuCycles++ {
			ppu.TickScanline()
		}

		if cpuCycles > FRAME_CYCLES {
			cpuCycles -= FRAME_CYCLES

			tv.ShowPixels()
			tv.UpdateInputState(controller)

			if delta := FRAME_TIME - float64(time.Now().UnixNano()-t0); delta > 0 {
				time.Sleep(time.Duration(delta))
      }

			t0 = time.Now().UnixNano()
		}
	}
}
