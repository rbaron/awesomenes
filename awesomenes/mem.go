package awesomenes

// In the future, maybe make this an interface, since
// we'll need memory mapped IO too
type Memory []uint8

func (m Memory) Read8(p uint16) uint8 {
  return m[p]
}

func (m Memory) Write8(p uint16, v uint8) {
  m[p] = v
}

