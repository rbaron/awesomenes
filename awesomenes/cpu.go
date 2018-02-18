package awesomenes

import (
  "fmt"
)

const (
  StatusFlagC = iota
  StatusFlagZ
  StatusFlagI
  StatusFlagD
  StatusFlagB
  StatusFlagV
  StatusFlagN
)

const (
  MemStackBase = 0x10ff
)

type registers struct {
  PC uint16
  SP uint8
  A  uint8
  X  uint8
  Y  uint8
  P  uint8
}

type CPU struct {
  regs *registers
  mem  AddrSpace
}

func makeCPU(addrSpace AddrSpace) *CPU {
  return &CPU{
    // Top of the stack
    regs:  &registers{
      SP: 0x0,
    },
    mem: addrSpace,
  }
}

func (cpu *CPU) Exec(rom *Rom) {
  opcode := cpu.mem.Read8(cpu.regs.PC)
  instr := instrTable[opcode]
  fmt.Println("Hello", instr)
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
  return c.mem.Read8(c.stackPos())
}

func (c *CPU) Pop16() uint16 {
  lsb := uint16(c.Pop8())
  msb := uint16(c.Pop8())
  return (msb << 8) + lsb
}

func (c *CPU) stackPos() uint16 {
  return MemStackBase - uint16(c.regs.SP)
}

func (c *CPU) getFlag(flag uint8) bool {
  return (c.regs.P & (0x1 << flag)) != 0
}

func (c *CPU) setFlag(flag uint8) {
  c.regs.P |= (0x1 << flag)
}

func (c *CPU) resetFlag(flag uint8) {
  c.regs.P &= ^(0x1 << flag)
}

func (c *CPU) setOrReset(flag uint8, cond bool) {
  if cond {
    c.setFlag(flag)
  } else {
    c.resetFlag(flag)
  }
}
