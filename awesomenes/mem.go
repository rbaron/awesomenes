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

func (m Memory) Read16(p uint16) uint16 {
  lo := uint16(m[p])
  hi := uint16(m[p+1])
  return (hi << 8) + lo
}

func (m Memory) Write16(p uint16, v uint16) {
  m[p]   = uint8(v & 0xff)
  m[p+1] = uint8(v >> 8)
}
