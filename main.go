package main

import (
  "github.com/rbaron/awesomenes/awesomenes"
  "os"
  "bufio"
  "time"
)

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func main() {
  tv := awesomenes.MakeTV()

  rom := awesomenes.ReadROM("smb.nes")
  ppu := awesomenes.NewPPU(nil, rom)
  ppu.Reset()

  cpuAddrSpace := awesomenes.MakeCPUAddrSpace(rom, ppu)

  cpu := awesomenes.MakeCPU(cpuAddrSpace)

  // Back ref to cpu from ppu so we can trigger NMIs
  ppu.CPU = cpu
  ppu.TV = tv

  cpu.PowerUp()

  // cpu.Run()
  frameCycles := 29781
  var cpuCycles, ppuCycles int

  var i = 0
  for {
    newCycles := cpu.Run()
    i++

    cpuCycles += newCycles

    for ppuCycles = 0; ppuCycles < 3 * newCycles; ppuCycles++ {
      ppu.Step()
    }

    if cpuCycles > frameCycles {
      cpuCycles = 0
      tv.ShowPixels()
      time.Sleep(100000)
    }
  }
  return
}
