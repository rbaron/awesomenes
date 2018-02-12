package awesomenes

import (
  "fmt"
)

const (
  StatusC = iota
  StatusZ
  StatusI
  StatusD
  StatusB
  StatusV
  StatusN
)

type registers struct {
  PC uint16
  SP uint8
  A  uint8
  X  uint8
  Y  uint8
  P  uint8
}

const (
  MemStackBase = 0x10ff
)

type CPU struct {
  regs *registers
  mem  Memory
}

func makeCPU() *CPU {
  return &CPU{
    // Top of the stack
    regs:  &registers{
      SP: 0x0,
    },
    mem: make(Memory, 0x10000),
  }
}

func (c *CPU) String() string {
  return fmt.Sprintf("<CPU Regs: %v>", c.regs)
}

func (c *CPU) Push8(v uint8) {
  c.mem.Write8(c.stackPos(), v)
  c.regs.SP++
}

func (c *CPU) Push16(v uint16) {
  c.Push8(uint8(v >> 8))
  c.Push8(uint8(v & 0xff))
}

func (c *CPU) Pop8() uint8 {
  c.regs.SP--
  v := c.mem.Read8(c.stackPos())
  return v
}

func (c *CPU) Pop16() uint16 {
  lsb := uint16(c.Pop8())
  msb := uint16(c.Pop8())
  return (msb << 8) + lsb
}

func (c *CPU) stackPos() uint16 {
  return MemStackBase - uint16(c.regs.SP)
}
