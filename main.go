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
  FRAME_TIME   = 1.0 / (float64(FRAME_RATE) * 1e9)
)

func main() {
  tv := awesomenes.MakeTV()

  if len(os.Args) != 2 {
    fmt.Println("Usage: awesomenes ROM_PATH")
    os.Exit(2)
  }

  rom := awesomenes.ReadROM(os.Args[1])
  ppu := awesomenes.NewPPU(nil, rom)
  ppu.Reset()

  controller := awesomenes.MakeController()

  cpuAddrSpace := awesomenes.MakeCPUAddrSpace(rom, ppu, controller)

  cpu := awesomenes.MakeCPU(cpuAddrSpace)

  // Back ref to cpu and tv from ppu so we can trigger NMIs and update frames
  ppu.CPU = cpu
  ppu.TV  = tv

  cpu.PowerUp()

  var cpuCycles, ppuCycles int

  t0 := time.Now().UnixNano()

  for {

    newCycles := cpu.Run()

    cpuCycles += newCycles

    for ppuCycles = 0; ppuCycles < 3 * newCycles; ppuCycles++ {
      ppu.TickScanline()
    }

    if cpuCycles > FRAME_CYCLES {
      cpuCycles -= FRAME_CYCLES

      tv.ShowPixels()

      tv.UpdateInputState(controller)

      delta := FRAME_TIME - float64(time.Now().UnixNano() - t0)

      if delta > 0 {
        time.Sleep(time.Duration(delta))
      }
      t0 = time.Now().UnixNano()
    }
  }
  return
}
