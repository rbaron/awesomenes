package awesomenes

import "log"

func addrSetter(v uint8, bitN uint8, ifNotSet uint16, ifSet uint16) uint16 {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}

func boolSetter(v uint8, bitN uint8, ifNotSet bool, ifSet bool) bool {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}

/*
  PPUCTRL
*/
const (
  SPRITE_SIZE_8    = false
  SPRITE_SIZE_16   = true
  MS_READ_EXT      = false
  MS_WRITE_EXT     = true
)

type PPUCTRL struct {
  NameTableAddr        uint16
  VRAMReadIncrement    uint16
  // Addr for 8x8 sprites only (ignored for 16x16)
  SpritePatTableAddr   uint16
  BgTableAddr          uint16
  SpriteSize           bool
  MasterSlave          bool
  NMIonVBlank          bool
}

func (ctrl *PPUCTRL) Set(v uint8) {
  switch v & 0x3 {
    case 0x0:
      ctrl.NameTableAddr = 0x2000
    case 0x1:
      ctrl.NameTableAddr = 0x2400
    case 0x2:
      ctrl.NameTableAddr = 0x2800
    case 0x3:
      ctrl.NameTableAddr = 0x2c00
  }

  ctrl.VRAMReadIncrement  = addrSetter(v, 2, 0x0001, 0x0020)
  ctrl.SpritePatTableAddr = addrSetter(v, 3, 0x0000, 0x1000)
  ctrl.BgTableAddr        = addrSetter(v, 4, 0x0000, 0x1000)
  ctrl.SpriteSize         = boolSetter(v, 5, SPRITE_SIZE_8, SPRITE_SIZE_16)
  ctrl.MasterSlave        = boolSetter(v, 6, MS_READ_EXT, MS_WRITE_EXT)
  ctrl.NMIonVBlank        = boolSetter(v, 7, false, true)
}

/*
  PPUMASK
*/

type PPUMASK struct {
  Greyscale       bool
  ShowBgLeft      bool
  ShowSpritesLeft bool
  showBg          bool
  showSprites     bool
  emphasisRed     bool
  emphasisGreen   bool
  emphasisBlue    bool
}

func (mask *PPUMASK) Set(v uint8) {
  mask.Greyscale       = boolSetter(v, 0, false, true)
  mask.ShowBgLeft      = boolSetter(v, 1, false, true)
  mask.ShowSpritesLeft = boolSetter(v, 2, false, true)
  mask.showBg          = boolSetter(v, 3, false, true)
  mask.showSprites     = boolSetter(v, 4, false, true)
  mask.emphasisRed     = boolSetter(v, 5, false, true)
  mask.emphasisGreen   = boolSetter(v, 6, false, true)
  mask.emphasisBlue    = boolSetter(v, 0, false, true)
}

type PPUSTATUS struct {
  SpriteOverflow bool
  Sprite0Hit     bool
  VBlankStarted  bool

  // So we can simulate a dirty bus when reading CTRL
  LastWrite uint8
}

func (status *PPUSTATUS) Get() (result uint8) {
  if status.SpriteOverflow {
    result |= 0x1 << 5
  }
  if status.Sprite0Hit {
    result |= 0x1 << 6
  }
  if status.VBlankStarted {
    result |= 0x1 << 7
  }

  result |= status.LastWrite & 0x1f
  return
}

type PPUADDR struct {
  addr    uint16
  writeLo bool
}

func (addr *PPUADDR) Write(v uint8) {
  if addr.writeLo {
    addr.addr |= uint16(v)
    addr.writeLo = false
  } else {
    addr.addr = uint16(v) << 8
    addr.writeLo = true
  }
}

type PPUSCROLL struct {
  X uint8
  Y uint8
  WriteY bool
}

func (scrl *PPUSCROLL) Write(v uint8) {
  if scrl.WriteY {
    scrl.Y = v
    scrl.WriteY = false
  } else {
    scrl.X = v
    scrl.WriteY = true
  }
}

type PPU struct {
  CTRL    *PPUCTRL
  MASK    *PPUMASK
  STATUS  *PPUSTATUS
  SCRL    *PPUSCROLL
  ADDR    *PPUADDR

  // This is usually mapped to be the chartridge ram!
  // On mapper 0, accessing 0-0x2000 on the PPU actually
  // accesses the cartridge's CHR-RAM/ROM
  //This is the "pattern table"
  PatternTableData Memory
  NametableData    Memory
  PaletteData      Memory


  OAMADDR uint8
  OAMData [256]uint8
}

func MakePPU(chrROM Memory) *PPU {
  return &PPU{
    CTRL:   &PPUCTRL{},
    MASK:   &PPUMASK{},
    STATUS: &PPUSTATUS{},

    PatternTableData: chrROM,
    NametableData:    make(Memory, 0x0800),
    PaletteData:      make(Memory, 0x0020),
  }
}

//OAMDATA
func (ppu *PPU) WriteOAMData(v uint8) {
  ppu.OAMData[ppu.OAMADDR] = v
  ppu.OAMADDR += 1
}

func (ppu *PPU) ReadOAMData() uint8 {
  return ppu.OAMData[ppu.OAMADDR]
}

//PPUDATA
func (ppu *PPU) WriteData(v uint8) {
  ppu.Write8(ppu.ADDR.addr, v)
  ppu.ADDR.addr += ppu.CTRL.VRAMReadIncrement
}

func (ppu *PPU) ReadData() uint8 {
  val := ppu.Read8(ppu.ADDR.addr)
  ppu.ADDR.addr += ppu.CTRL.VRAMReadIncrement
  return val
}

func (ppu *PPU) OMADMA(data []uint8) {
  copy(ppu.OAMData[:], data)
}

func (ppu *PPU) Write8 (addr uint16, v uint8) {
  switch {
    // Pattern tables - for now hard mapped to CHRROM
    case addr >= 0x0000 && addr < 0x2000:
      ppu.PatternTableData.Write8(addr, v)

    case addr >= 0x2000 && addr < 0x3f00:
      ppu.NametableData.Write8(getMirroedAddr(addr), v)

    case addr >= 0x3f00 && addr < 0x4000:
      ppu.PaletteData.Write8(addr, v)

    default:
      log.Fatalf("Invalid write to PPU at %x", addr)
  }
}

func (ppu *PPU) Read8(addr uint16) uint8 {
  switch {
    // Pattern tables - for now hard mapped to CHRROM
    case (addr >= 0x0000 && addr < 0x2000):
      return ppu.PatternTableData.Read8(addr)

    case addr >= 0x2000 && addr < 0x3f00:
      return ppu.NametableData.Read8(getMirroedAddr(addr))

    case addr >= 0x3f00 && addr < 0x4000:
      return ppu.PaletteData.Read8(addr)

    default:
      log.Fatalf("Invalid read from PPU at %x", addr)
      return 0
  }
}

// Hard coded vertical mirror for now
func getMirroedAddr(addr uint16) uint16 {
  v := (addr - 0x2000) % 0x800
  return 0x2000 + v
}
