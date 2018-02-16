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
      return uint16(cpu.mem.Read8(uint16(m + x)))

    case AddrModeIndirectY:
      m := cpu.mem.Read8(cpu.regs.PC + 1)
      y := cpu.regs.Y
      return uint16(cpu.mem.Read8(uint16(m)) + y)

    case AddrModeZeroPage:
      return uint16(cpu.mem.Read8(cpu.regs.PC + 1))

    // TODO: difference from AddrModeXIndirect
    case AddrModeZeroX:
      return uint16(cpu.mem.Read8(cpu.regs.PC + 1) + cpu.regs.X)

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

func eor(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.regs.A = cpu.regs.A ^ cpu.mem.Read8(addr)
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

// Pull processor state
func plp(cpu *CPU, addrMode addressingMode) {
  cpu.regs.P = cpu.Pop8()
}

// Clear carry
func clc(cpu *CPU, addrMode addressingMode) {
  cpu.resetFlag(StatusFlagC)
}

func cld(cpu *CPU, addrMode addressingMode) {
  cpu.resetFlag(StatusFlagD)
}

func cli(cpu *CPU, addrMode addressingMode) {
  cpu.resetFlag(StatusFlagI)
}

func clv(cpu *CPU, addrMode addressingMode) {
  cpu.resetFlag(StatusFlagV)
}

// Jump
func jmp(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.regs.PC = addr
}

// Jump to subroutine
func jsr(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.Push16(cpu.regs.PC - 1)
  cpu.regs.PC = addr
}

func and(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  cpu.regs.A = cpu.regs.A & cpu.mem.Read8(addr)
  cpu.setOrReset(StatusFlagN, cpu.regs.A & 0x80 != 0)
  cpu.setOrReset(StatusFlagZ, cpu.regs.A == 0)
}


func bit(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  v := cpu.mem.Read8(addr)

  cpu.setOrReset(StatusFlagZ, cpu.regs.A & v == 0)
  cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
  cpu.setOrReset(StatusFlagV, (v >> 6) & 0x1 == 0x1)
}

// Rotate left
func rol(cpu *CPU, addrMode addressingMode) {
  inner := func (v uint8) uint8 {
    v = (v << 1) | (v >> 7)
    cpu.setOrReset(StatusFlagZ, v == 0)
    cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
    return v
  }

  if addrMode == AddrModeAccumulator {
    cpu.regs.A = inner(cpu.regs.A)
  } else {
    addr := calculateAddr(cpu, addrMode)
    cpu.mem.Write8(addr, inner(cpu.mem.Read8(addr)))
  }
}

// Rotate right
func ror(cpu *CPU, addrMode addressingMode) {
  inner := func (v uint8) uint8 {
    v = (v >> 1) | ((v & 0x1) << 7)
    cpu.setOrReset(StatusFlagZ, v == 0)
    cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
    return v
  }

  if addrMode == AddrModeAccumulator {
    cpu.regs.A = inner(cpu.regs.A)
  } else {
    addr := calculateAddr(cpu, addrMode)
    cpu.mem.Write8(addr, inner(cpu.mem.Read8(addr)))
  }
}

// Branch if negative
func bmi(cpu *CPU, addrMode addressingMode) {
  if cpu.getFlag(StatusFlagN) {
    cpu.regs.PC = calculateAddr(cpu, addrMode)
  }
}

// Branch if not equal
func bne(cpu *CPU, addrMode addressingMode) {
  if !cpu.getFlag(StatusFlagZ) {
    cpu.regs.PC = calculateAddr(cpu, addrMode)
  }
}

// Branch if positive
func bpl(cpu *CPU, addrMode addressingMode) {
  if !cpu.getFlag(StatusFlagN) {
    cpu.regs.PC = calculateAddr(cpu, addrMode)
  }
}

// Branch if overflow clear
func bvc(cpu *CPU, addrMode addressingMode) {
  if !cpu.getFlag(StatusFlagV) {
    cpu.regs.PC = calculateAddr(cpu, addrMode)
  }
}

// Branch if overflow set
func bvs(cpu *CPU, addrMode addressingMode) {
  if cpu.getFlag(StatusFlagV) {
    cpu.regs.PC = calculateAddr(cpu, addrMode)
  }
}

func _comp(v uint8, cpu *CPU) {
  cpu.setOrReset(StatusFlagZ, v == 0)
  cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
  cpu.setOrReset(StatusFlagC, v >= 0)
}

// Compare with A
func cmp(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  m := cpu.mem.Read8(addr)
  v := cpu.regs.A - m
  _comp(v, cpu)
}

// Compare with X
func cpx(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  m := cpu.mem.Read8(addr)
  v := cpu.regs.X - m
  _comp(v, cpu)
}

// Compare with Y
func cpy(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  m := cpu.mem.Read8(addr)
  v := cpu.regs.Y - m
  _comp(v, cpu)
}

func dec(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  v := cpu.mem.Read8(addr) - 1
  cpu.mem.Write8(addr, v)
  cpu.setOrReset(StatusFlagZ, v == 0)
  cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
}

func dex(cpu *CPU, addrMode addressingMode) {
  cpu.regs.X -= 1
  cpu.setOrReset(StatusFlagZ, cpu.regs.X == 0)
  cpu.setOrReset(StatusFlagN, cpu.regs.X >> 7 == 0x1)
}

func dey(cpu *CPU, addrMode addressingMode) {
  cpu.regs.Y -= 1
  cpu.setOrReset(StatusFlagZ, cpu.regs.Y == 0)
  cpu.setOrReset(StatusFlagN, cpu.regs.Y >> 7 == 0x1)
}

func inc(cpu *CPU, addrMode addressingMode) {
  addr := calculateAddr(cpu, addrMode)
  v := cpu.mem.Read8(addr) + 1
  cpu.mem.Write8(addr, v)
  cpu.setOrReset(StatusFlagZ, v == 0)
  cpu.setOrReset(StatusFlagN, v >> 7 == 0x1)
}

func inx(cpu *CPU, addrMode addressingMode) {
  cpu.regs.X += 1
  cpu.setOrReset(StatusFlagZ, cpu.regs.X == 0)
  cpu.setOrReset(StatusFlagN, cpu.regs.X >> 7 == 0x1)
}

func iny(cpu *CPU, addrMode addressingMode) {
  cpu.regs.Y += 1
  cpu.setOrReset(StatusFlagZ, cpu.regs.Y == 0)
  cpu.setOrReset(StatusFlagN, cpu.regs.Y >> 7 == 0x1)
}

func lda(cpu *CPU, addrMode addressingMode) {
}
