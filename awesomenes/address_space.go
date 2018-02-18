package awesomenes

import (
  "log"
)

type AddrSpace interface {
  Read8(addr uint16) uint8
  Write8(addr uint16, v uint8)

  Read16(addr uint16) uint16
  Write16(addr uint16, v uint16)
}

type CPUAddrSpace struct {
  RAM Memory
  ROM *Rom

  // Placeholder until we have a proper PPU implementation
  PPU Memory

  //APU

  // Logger for tests
  Logger Memory

  // Mapper
  // http://tuxnes.sourceforge.net/nesmapper.txt
}

func makeCPUAddrSpace(rom *Rom) *CPUAddrSpace {
  return &CPUAddrSpace{
    RAM:    make(Memory, 0x800),
    ROM:    rom,
    Logger: make(Memory, 0x1000),
  }
}

//http://wiki.nesdev.com/w/index.php/CPU_memory_map
func (as *CPUAddrSpace) Read8(addr uint16) uint8 {

  switch {
    case addr < 0x2000:
      // 0x0800 - 0x1fff mirrors 0x0000 - 0x07ff three times
      return as.RAM.Read8(addr % 0x800)

    // PPU registers
    case addr < 0x4000:
      return as.PPU.Read8(0x2000 + addr % 8)

    default:
      log.Printf("Unhandled read from CPU mem space at %x", addr)
      return 0
  }
}

func (as *CPUAddrSpace) Write8(addr uint16, v uint8) {

  switch {
    case addr < 0x2000:
      // 0x0800 - 0x1fff mirrors 0x0000 - 0x07ff three times
      as.RAM.Write8(addr % 0x800, v)

    // PPU registers
    case addr < 0x4000:
      as.PPU.Write8(0x2000 + addr % 8, v)

    default:
      log.Printf("Unhandled write to CPU mem space at %x", addr)
  }
}

// Little-endian mem layout
func (as *CPUAddrSpace) Read16(addr uint16) uint16 {
  lo := uint16(as.Read8(addr))
  hi := uint16(as.Read8(addr + 1))
  return (hi << 8) + lo
}

func (as *CPUAddrSpace) Write16(addr uint16, v uint16) {
  as.Write8(addr, uint8(v & 0xff))
  as.Write8(addr + 1, uint8(v >> 8))
}
