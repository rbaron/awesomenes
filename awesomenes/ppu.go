package awesomenes

import (
  "log"
)

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
  // TODO set temp addr?
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

func (ppu *PPU) LowBGTileAddr() uint16 {
  return ppu.CTRL.BgTableAddr + uint16(ppu.NameTableLatch) * 16 + ppu.ADDR.FineY()
}

func (ppu *PPU) HighBGTileAddr() uint16 {
  return ppu.LowBGTileAddr() + 8
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
  VAddr        uint16  // v
  TAddr        uint16  // t
  WriteHi      bool    // w
  FineXScroll  uint8   // x
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) NameTableAddr() uint16 {
  return 0x2000 | (addr.VAddr & 0x0fff)
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) AttrTableAddr() uint16 {
  v := addr.VAddr
  return 0x23c0 | (v & 0x0c00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
}

func (addr *PPUADDR) FineY() uint16 {
  return (addr.VAddr >> 12) & 0x07
}


func (addr *PPUADDR) Write(v uint8) {
  if addr.WriteHi == false {
    addr.TAddr |= uint16(v) << 8
    addr.TAddr &= 0x7fff
    addr.WriteHi = true
  } else {
    addr.TAddr |= uint16(v)
    addr.VAddr = addr.TAddr
    addr.WriteHi = false
  }
}

func (addr *PPUADDR) SetOnCTRLWrite(v uint8) {
  addr.TAddr |= uint16(v & 0x03) << 10
}

func (addr *PPUADDR) SetOnSTATUSRead() {
  addr.WriteHi = false
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling
func (addr *PPUADDR) TransferX () {
  addr.VAddr = (addr.VAddr & 0xFBE0) | (addr.TAddr & 0x041F)
}

func (addr *PPUADDR) TransferY () {
  addr.VAddr = (addr.VAddr & 0x841F) | (addr.TAddr & 0x7BE0)
}

func (addr *PPUADDR) SetOnSCROLLWrite(v uint8) {
  if addr.WriteHi == false {
    addr.TAddr |= uint16(v >> 3)
    addr.FineXScroll = v & 0x3
    addr.WriteHi = true
  } else {
    addr.TAddr |= uint16(v & 0x03) << 12
    addr.TAddr |= uint16(v & 0xf8) << 2
    addr.WriteHi = false
  }
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling#Y_increment
func (addr *PPUADDR) IncrementFineY() {
  v := addr.VAddr
  var y uint16

  if ((v & 0x7000) != 0x7000) {
    v += 0x1000
  } else {
    //v &= ^0x7000
    v &= 0x8FFF
    y = (v & 0x03E0) >> 5
    if (y == 29) {
      y = 0
      v ^= 0x0800
    } else {
      if (y == 31) {
        y = 0
      } else {
        y += 1
      }
    }
  }
  addr.VAddr  = (v & 0xFC1F) | (y << 5)
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling#X_increment
func (addr *PPUADDR) IncrementCoarseX() {
  v := addr.VAddr

  if ((v & 0x001F) == 31) {
    v &= 0xFFE0
    v ^= 0x0400
  } else {
    v += 1
  }

  addr.VAddr = v
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

  CPU     *CPU

  tv      *TV

  // This is usually mapped to be the chartridge ram!
  // On mapper 0, accessing 0-0x2000 on the PPU actually
  // accesses the cartridge's CHR-RAM/ROM
  //This is the "pattern table"
  PatternTableData Memory
  NametableData    Memory
  PaletteData      Memory


  OAMADDR uint8
  OAMData [256]uint8

  /*
    Rendering
  */
  Scanline  int
  Dot       int

  // Latches that will be fetched during visible/pre scanlines
  // and then pushed into the
  NameTableLatch  uint8
  AttrTableLatch  uint8
  BgLatchLow      uint8
  BgLatchHigh     uint8

  // Background
  //VRAMAddr     uint16
  //VRAMAddrTemp uint16
  //FineXScroll  uint16
  //IsFirstWrite bool

  BgTileShiftLow    uint16
  BgTileShiftHigh   uint16

  // Low byte for bg that will be put in the shift reg

  BgPaletteShift1 uint8
  BgPaletteShift2 uint8

  // Sprite
  PrimaryOAMBuffer   []OAMSprite
  SecondaryOAMBuffer []OAMSprite

  Pixels             []byte
}

type OAMSprite struct {
  x    uint8
  y    uint8

  // Number of the tile in the pattern tables
  tile uint8

  // Links this sprite to a color palette
  attr uint8
}

func MakePPU(chrROM Memory, tv *TV) *PPU {
  return &PPU{
    ADDR:   &PPUADDR{},
    CTRL:   &PPUCTRL{},
    MASK:   &PPUMASK{},
    STATUS: &PPUSTATUS{},
    SCRL:   &PPUSCROLL{},

    tv:     tv,

    PatternTableData: chrROM,
    NametableData:    make(Memory, 0x0800),
    PaletteData:      make(Memory, 0x0020),

    PrimaryOAMBuffer:   make([]OAMSprite, 64),
    SecondaryOAMBuffer: make([]OAMSprite,  8),

    Pixels: make([]byte, 256 * 240),
  }
}

func (ppu *PPU) Run() {
  ppu.TickScanline()
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
  ppu.Write8(ppu.ADDR.VAddr, v)
  ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement
}

func (ppu *PPU) ReadData() uint8 {
  val := ppu.Read8(ppu.ADDR.VAddr)
  ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement
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
      ppu.PaletteData.Write8(addr % 0x20, v)

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
      return ppu.PaletteData.Read8(addr % 0x20)

    default:
      log.Fatalf("Invalid read from PPU at %x", addr)
      return 0
  }
}

// Hard coded vertical mirror for now
func getMirroedAddr(addr uint16) uint16 {
  return addr % 0x800
}
