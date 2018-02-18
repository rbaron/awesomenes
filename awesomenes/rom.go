package awesomenes

import (
  "io/ioutil"
)

//http://fms.komkon.org/EMUL8/NES.html
type Rom struct {
  header    *RomHeader
  rom       []uint8
  vrom      []uint8
}

type RomHeader struct {
  // 16kB each
  nROMBanks  uint8
  // 8kB each
  nVROMBanks uint8
}

func ReadROM(path string) *Rom {
  data, _ := ioutil.ReadFile(path)

  if string(data[:3]) != "NES" {
    panic("Invalid ROM file" + string(data[:3]))
  }

  header := &RomHeader{
    nROMBanks:  data[4],
    nVROMBanks: data[5],
  }

  rom := &Rom{
    header: header,
    rom:    data[16:(int(header.nROMBanks)*0x4000)],
    //vrom:   data[]
  }

  return rom
}
