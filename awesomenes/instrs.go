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
    name:     "PLP",
    opcode:   0x28,
    size:     1,
    cycles:   4,
    addrMode: AddrModeImplied,
    fn:       plp,
  },
  &instr{
    name:     "CLC",
    opcode:   0x18,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       clc,
  },
  &instr{
    name:     "CLD",
    opcode:   0xd8,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       cld,
  },
  &instr{
    name:     "CLI",
    opcode:   0x58,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       cli,
  },
  &instr{
    name:     "CLV",
    opcode:   0xb8,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       clv,
  },
  &instr{
    name:     "JSR",
    opcode:   0x20,
    size:     3,
    cycles:   6,
    addrMode: AddrModeAbs,
    fn:       jsr,
  },
  &instr{
    name:     "AND",
    opcode:   0x29,
    size:     2,
    cycles:   2,
    addrMode: AddrModeImmediate,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x25,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x35,
    size:     2,
    cycles:   4,
    addrMode: AddrModeZeroX,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x2d,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x3d,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsX,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x39,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsY,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x21,
    size:     2,
    cycles:   6,
    addrMode: AddrModeXIndirect,
    fn:       and,
  },
  &instr{
    name:     "AND",
    opcode:   0x31,
    size:     2,
    cycles:   5,
    addrMode: AddrModeIndirectY,
    fn:       and,
  },
  &instr{
    name:     "BIT",
    opcode:   0x24,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       bit,
  },
  &instr{
    name:     "BIT",
    opcode:   0x2c,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       bit,
  },
  &instr{
    name:     "ROL",
    opcode:   0x2a,
    size:     1,
    cycles:   2,
    addrMode: AddrModeAccumulator,
    fn:       rol,
  },
  &instr{
    name:     "ROL",
    opcode:   0x26,
    size:     2,
    cycles:   5,
    addrMode: AddrModeZeroPage,
    fn:       rol,
  },
  &instr{
    name:     "ROL",
    opcode:   0x36,
    size:     2,
    cycles:   6,
    addrMode: AddrModeZeroX,
    fn:       rol,
  },
  &instr{
    name:     "ROL",
    opcode:   0x2e,
    size:     3,
    cycles:   6,
    addrMode: AddrModeAbs,
    fn:       rol,
  },
  &instr{
    name:     "ROL",
    opcode:   0x3e,
    size:     3,
    cycles:   7,
    addrMode: AddrModeAbsX,
    fn:       rol,
  },
  &instr{
    name:     "ROR",
    opcode:   0x6a,
    size:     1,
    cycles:   2,
    addrMode: AddrModeAccumulator,
    fn:       ror,
  },
  &instr{
    name:     "ROR",
    opcode:   0x66,
    size:     2,
    cycles:   5,
    addrMode: AddrModeZeroPage,
    fn:       ror,
  },
  &instr{
    name:     "ROR",
    opcode:   0x76,
    size:     2,
    cycles:   6,
    addrMode: AddrModeZeroX,
    fn:       ror,
  },
  &instr{
    name:     "ROR",
    opcode:   0x6e,
    size:     3,
    cycles:   6,
    addrMode: AddrModeAbs,
    fn:       ror,
  },
  &instr{
    name:     "ROR",
    opcode:   0x7e,
    size:     3,
    cycles:   7,
    addrMode: AddrModeAbsX,
    fn:       ror,
  },
  &instr{
    name:     "BMI",
    opcode:   0x30,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bmi,
  },
  &instr{
    name:     "BNE",
    opcode:   0xd0,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bne,
  },
  &instr{
    name:     "BPL",
    opcode:   0x10,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bpl,
  },
  &instr{
    name:     "BVC",
    opcode:   0x50,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bvc,
  },
  &instr{
    name:     "BVS",
    opcode:   0x70,
    size:     2,
    cycles:   2,
    addrMode: AddrModeRelative,
    fn:       bvc,
  },
  &instr{
    name:     "CMP",
    opcode:   0xc9,
    size:     2,
    cycles:   2,
    addrMode: AddrModeImmediate,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xc5,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xd5,
    size:     2,
    cycles:   4,
    addrMode: AddrModeZeroX,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xcd,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xdd,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsX,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xd9,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbsY,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xc1,
    size:     2,
    cycles:   6,
    addrMode: AddrModeXIndirect,
    fn:       cmp,
  },
  &instr{
    name:     "CMP",
    opcode:   0xd1,
    size:     2,
    cycles:   5,
    addrMode: AddrModeIndirectY,
    fn:       cmp,
  },
  &instr{
    name:     "CPX",
    opcode:   0xe0,
    size:     2,
    cycles:   2,
    addrMode: AddrModeImmediate,
    fn:       cpx,
  },
  &instr{
    name:     "CPX",
    opcode:   0xe4,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       cpx,
  },
  &instr{
    name:     "CPX",
    opcode:   0xec,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       cpx,
  },
  &instr{
    name:     "CPY",
    opcode:   0xc0,
    size:     2,
    cycles:   2,
    addrMode: AddrModeImmediate,
    fn:       cpy,
  },
  &instr{
    name:     "CPY",
    opcode:   0xc4,
    size:     2,
    cycles:   3,
    addrMode: AddrModeZeroPage,
    fn:       cpy,
  },
  &instr{
    name:     "CPY",
    opcode:   0xcc,
    size:     3,
    cycles:   4,
    addrMode: AddrModeAbs,
    fn:       cpy,
  },
  &instr{
    name:     "DEC",
    opcode:   0xc6,
    size:     2,
    cycles:   5,
    addrMode: AddrModeZeroPage,
    fn:       dec,
  },
  &instr{
    name:     "DEC",
    opcode:   0xd6,
    size:     2,
    cycles:   6,
    addrMode: AddrModeZeroX,
    fn:       dec,
  },
  &instr{
    name:     "DEC",
    opcode:   0xce,
    size:     3,
    cycles:   6,
    addrMode: AddrModeAbs,
    fn:       dec,
  },
  &instr{
    name:     "DEC",
    opcode:   0xde,
    size:     3,
    cycles:   7,
    addrMode: AddrModeAbsX,
    fn:       dec,
  },
  &instr{
    name:     "DEX",
    opcode:   0xca,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       dex,
  },
  &instr{
    name:     "DEY",
    opcode:   0x88,
    size:     1,
    cycles:   2,
    addrMode: AddrModeImplied,
    fn:       dey,
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
