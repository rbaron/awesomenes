package awesomenes

import (
  "testing"
)

func TestPushPop(t *testing.T) {
  cpu := makeCPU()

  cpu.Push8(0xde)
  cpu.Push8(0xad)
  cpu.Push16(0xbeaf)

  if v := cpu.Pop16(); v != 0xbeaf {
    t.Fatalf("Wrong value for Pop16: %x", v)
  }

  if v := cpu.Pop8(); v != 0xad {
    t.Fatalf("Wrong value for Pop8: %x", v)
  }

  if v := cpu.Pop8(); v != 0xde {
    t.Fatalf("Wrong value for Pop8: %x", v)
  }

}
