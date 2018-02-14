package awesomenes

import (
  "log"
)

type instr struct {
  name     string
  opcode   uint8
  size     uint8
  cycles   uint8
  addrMode addressingMode
  fn       func(cpu *CPU, addrMode addressingMode)
}

type addressingMode uint8

const (
  AddrModeAbs          addressingMode = iota
  AddrModeAbsX
  AddrModeAbsY
  AddrModeAccumulator
  AddrModeImmediate
  AddrModeImplied
  AddrModeIndirect
  AddrModeXIndirect
  AddrModeIndirectY
  AddrModeRelative
  AddrModeZeroPage
  AddrModeZeroX
  AddrModeZeroY
)

var instrs = []*instr {
  &instr{
    name:     "BRK",
    opcode:   0x00,
    size:     1,
    cycles:   7,
    addrMode: AddrModeImplied,
    fn:       brk,
  },
  &instr{
    name:     "ORA",
    opcode:   0x01,
    size:     2,
    cycles:   6,
    addrMode: AddrModeXIndirect,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x05,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x09,
    size:     2,
    cycles:   2,
    addrMode: AddrModeImmediate,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x0d,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x11,
    size:     2,
    cycles:   5,
    addrMode: AddrModeIndirectY,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x15,
    size:     2,
    cycles:   4,
    addrMode: AddrModeZeroX,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x19,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsY,
    fn:       ora,
  },
  &instr{
    name:     "ORA",
    opcode:   0x1d,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsX,
    fn:       ora,
  },
  &instr{
    name:     "ASL",
    opcode:   0x0a,
    size:     1,
    cycles:   2,
    addrMode: AddrModeAccumulator,
    fn:       asl,
  },
  &instr{
    name:     "ASL",
    opcode:   0x06,
    size:     2,
    cycles:   5,
    addrMode: AddrModeZeroPage,
    fn:       asl,
  },
  &instr{
    name:     "ASL",
    opcode:   0x16,
    size:     2,
    cycles:   6,
    addrMode: AddrModeZeroX,
    fn:       asl,
  },
  &instr{
    name:     "ASL",
    opcode:   0x0e,
    size:     3,
    cycles:   6,
    addrMode: AddrModeAbs,
    fn:       asl,
  },
  &instr{
    name:     "ASL",
    opcode:   0x1e,
    size:     3,
    cycles:   7,
    addrMode: AddrModeAbsX,
    fn:       asl,
  },
  &instr{
    name:     "PHP",
    opcode:   0x08,
    size:     1,
    cycles:   3,
    addrMode: AddrModeImplied,
    fn:       php,
  },
  &instr{
    name:     "BPL",
    opcode:   0x10,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bpl,
  },
}

func calculateAddr(cpu *CPU, addrMode addressingMode) uint16 {
  switch addrMode {
    // TODO: does endianess matter with the abs modes? Currently using little-endian.
    case AddrModeAbs:
      return cpu.mem.Read16(cpu.regs.PC + 1)

    case AddrModeAbsX:
      return cpu.mem.Read16(cpu.regs.PC + 1) + uint16(cpu.regs.X)

    case AddrModeAbsY:
      return cpu.mem.Read16(cpu.regs.PC + 1) + uint16(cpu.regs.Y)

    case AddrModeAccumulator:
      log.Fatalf("It makes no sense to calculate addresses in accumulator addressing mode")
      return 0xffff

    case AddrModeImmediate:
      return cpu.regs.PC + 1

    case AddrModeIndirectY:
      return uint16(cpu.regs.PC + 1) + uint16(cpu.regs.Y)

    case AddrModeRelative:
      // Treat operand as signed int8. Sure it uses two's complement?
      m := cpu.mem.Read8(cpu.regs.PC + 1)
      if (m >> 7) == 0x1 {
        return cpu.regs.PC + 2 + uint16(m) - 0x100
      } else {
        return cpu.regs.PC + 2 + uint16(m)
      }

    case AddrModeXIndirect:
      m := cpu.mem.Read8(cpu.regs.PC + 1)
      x := cpu.regs.X
      return uint16(m + x)

    case AddrModeZeroPage:
      return uint16(cpu.regs.PC + 1)

    case AddrModeZeroX:
      return uint16(uint8(cpu.regs.PC) + 1 + cpu.regs.X)

    default:
      log.Fatalf("Invalid addressing mode")
      return 0xffff
  }
}

// Break
func brk(cpu *CPU, addrMode addressingMode) {
  cpu.Push16(cpu.regs.PC)
  cpu.Push8(cpu.regs.P)
  cpu.setFlag(StatusFlagB)
  cpu.regs.PC = cpu.mem.Read16(0xfffe)
}

// Logical inclusive OR
func ora(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.regs.A = cpu.regs.A | cpu.mem.Read8(addr)
  cpu.setOrReset(StatusFlagN, cpu.regs.A & 0x80 != 0)
  cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
}

// Arithmetic shift left
func asl(cpu *CPU, addrMode addressingMode) {
  shiftL := func (v uint8) uint8 {
    cpu.setOrReset(StatusFlagC, v & 0x80 != 0)
    v = v << 1
    cpu.setOrReset(StatusFlagZ, v == 0)
    return v
  }

  if addrMode == AddrModeAccumulator {
    cpu.regs.A = shiftL(cpu.regs.A)
  } else {
    addr := calculateAddr(cpu, addrMode)
    cpu.mem.Write8(addr, shiftL(cpu.mem.Read8(addr)))
  }
}

// Push processor state
func php(cpu *CPU, addrMode addressingMode) {
  cpu.Push8(cpu.regs.P)
}

// Branch if positive
func bpl(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.regs.PC = addr
}
