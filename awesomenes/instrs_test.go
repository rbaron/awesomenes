package awesomenes

import (
  "testing"
)

func TestBRK(t *testing.T) {
  cpu := makeCPU()
  cpu.mem.Write8(0xfffe, 0xad)
  cpu.mem.Write8(0xffff, 0xde)

  brk(cpu, AddrModeImplied)

  if cpu.regs.P & (0x1 << StatusFlagB) == 0 {
    t.Fatalf("Wrong value 0 for status bit B")
  }
  if cpu.regs.PC != 0xdead {
    t.Fatalf("Wrong value for PC register")
  }
}

func TestORA(t *testing.T) {
  cpu := makeCPU()
  cpu.mem.Write8(cpu.regs.PC + 1, 0xad)
  cpu.regs.A = 0x4a
  cpu.regs.X = 0x01

  // AddrModeXIndirect will read from mem[PC + 1] | X
  cpu.mem.Write8(0xae, 0x8d)

  ora(cpu, AddrModeXIndirect)

  if cpu.regs.A != uint8(0x4a | 0x8d) {
    t.Fatalf("Wrong value for reg A: %x", cpu.regs.A)
  }
  if cpu.regs.P & (0x1 << StatusFlagN) == 0 {
    t.Fatalf("Flag N shouldve been set")
  }
  if cpu.regs.P & (0x1 << StatusFlagZ) != 0 {
    t.Fatalf("Flag Z should not have been set")
  }
}

func TestASL(t *testing.T) {
  cpu := makeCPU()
  cpu.regs.A = 0x8a

  asl(cpu, AddrModeAccumulator)

  if cpu.regs.A != uint8((0x8a << 1) & 0xff) {
    t.Fatalf("Wrong value for reg A: %x", cpu.regs.A)
  }
  if cpu.regs.P & (0x1 << StatusFlagC) == 0 {
    t.Fatalf("Flag C shouldve been set")
  }
  if cpu.regs.P & (0x1 << StatusFlagZ) != 0 {
    t.Fatalf("Flag Z should not have been set")
  }
}

func TestBPL(t *testing.T) {
  cpu := makeCPU()

  cpu.regs.PC = 0x0004

  // 0xfa = signed -6
  cpu.mem.Write8(cpu.regs.PC + 1, 0xfa)

  bpl(cpu, AddrModeRelative)

  // PC should be at 0x0000 (relative jump of -6 plus 2 for the
  // BPL instruct itself)
  if cpu.regs.PC != 0x0000 {
    t.Fatalf("Wrong value for reg PC: %x", cpu.regs.PC)
  }

  cpu.regs.PC = 0x0004

  // 0xfa = signed +6
  cpu.mem.Write8(cpu.regs.PC + 1, 0x06)

  bpl(cpu, AddrModeRelative)

  // PC should be at 0x000c (relative jump of +6 plus 2 for the
  // BPL instruct itself)
  if cpu.regs.PC != 0x000c {
    t.Fatalf("Wrong value for reg PC: %x", cpu.regs.PC)
  }
}
