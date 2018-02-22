package awesomenes

import (
  "io/ioutil"
)

//http://fms.komkon.org/EMUL8/NES.html
type Rom struct {
  Header    *RomHeader
  PRGROM    Memory
  CHRROM    Memory
  PRGRAM    Memory
}

type RomHeader struct {
  MapperN        uint8
  // 16kB each
  NPRGROMBanks   uint8
  // 8kB each
  NCHRROMBanks   uint8

  HasTrainer     bool

  VerticalMirror bool
}

func ReadROM(path string) *Rom {
  data, _ := ioutil.ReadFile(path)

  if string(data[:3]) != "NES" {
    panic("Invalid ROM file" + string(data[:3]))
  }

  header := &RomHeader{
    MapperN:        (data[6] >> 4) | (data[7] & 0xf0),
    NPRGROMBanks:   data[4],
    NCHRROMBanks:   data[5],
    HasTrainer:     (data[6] & (0x1 << 2)) > 0,
    VerticalMirror: data[6] & 0x1 == 0x1,
  }

  if header.MapperN != 0 {
    panic("Only mapper type 0 is supported so far: " + string(header.MapperN));
  }

  if header.NPRGROMBanks != 2 {
    panic("Only 2 prg rom banks supported")
  }

  if header.NCHRROMBanks != 1 {
    panic("Only 1 chr rom banks supported")
  }

  var (
    prgBeginning uint16 = 16
    prgEnd       uint16 = 16 + uint16(header.NPRGROMBanks) * 0x4000
  )

  if header.HasTrainer {
    prgBeginning += 512
    prgEnd       += 512
  }

  var (
    chrBeginning uint16 = prgEnd
    chrEnd       uint16 = prgEnd + uint16(header.NCHRROMBanks) * 0x2000
  )

  rom := &Rom{
    Header: header,
    PRGROM: data[prgBeginning:prgEnd],
    CHRROM: data[chrBeginning:chrEnd],
    PRGRAM: make(Memory, 0x2000),
  }

  //rom.ROM.Dump(0, 256)

  return rom
}
