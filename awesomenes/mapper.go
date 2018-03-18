package awesomenes

type Mapper interface {
	Read8(addr uint16) uint8
	Write8(addr uint16, v uint8)
}
