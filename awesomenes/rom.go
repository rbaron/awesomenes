package awesomenes

import (
  "io/ioutil"
)

//http://fms.komkon.org/EMUL8/NES.html
type Rom struct {
  header    *RomHeader
  ROM       Memory
  VROM      Memory
  SRAM      Memory
}

type RomHeader struct {
  mapperN    uint8
  // 16kB each
  nROMBanks  uint8
  // 8kB each
  nVROMBanks uint8

  hasTrainer bool
}

func ReadROM(path string) *Rom {
  data, _ := ioutil.ReadFile(path)

  if string(data[:3]) != "NES" {
    panic("Invalid ROM file" + string(data[:3]))
  }

  header := &RomHeader{
    mapperN:    (data[6] >> 4) | (data[7] & 0xf0),
    nROMBanks:  data[4],
    nVROMBanks: data[5],
    hasTrainer: (data[6] & (0x1 << 2)) > 0,
  }

  if header.mapperN != 0 {
    panic("Only mapper type 0 is supported so far: " + string(header.mapperN));
  }

  if header.nROMBanks != 2 {
    panic("Only 2 rom banks supported")
  }

  var (
    romBeginningAddr uint16 = 16
    romEndAddr       uint16 = 16 + uint16(header.nROMBanks) * 0x4000
  )

  if header.hasTrainer {
    romBeginningAddr += 512
    romEndAddr       += 512
  }

  rom := &Rom{
    header: header,
    ROM:    data[romBeginningAddr:romEndAddr],
    //vrom:   data[]

    // Always 2kB of RAM for now
    SRAM:   make(Memory, 0x800),
  }

  rom.ROM.Dump(0, 256)

  return rom
}
