package awesomenes

/*
import (
  "log"
  "image"
)

func addrSetter(v uint8, bitN uint8, ifNotSet uint16, ifSet uint16) uint16 {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}

func boolSetter(v uint8, bitN uint8, ifNotSet bool, ifSet bool) bool {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}


// PPUCTRL

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
  //log.Printf("WROTE CONTROL %b", v)
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

// PPUMASK

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
  //log.Printf("WROTE MASK %b", v)
  mask.Greyscale       = boolSetter(v, 0, false, true)
  mask.ShowBgLeft      = boolSetter(v, 1, false, true)
  mask.ShowSpritesLeft = boolSetter(v, 2, false, true)
  mask.showBg          = boolSetter(v, 3, false, true)
  mask.showSprites     = boolSetter(v, 4, false, true)
  mask.emphasisRed     = boolSetter(v, 5, false, true)
  mask.emphasisGreen   = boolSetter(v, 6, false, true)
  mask.emphasisBlue    = boolSetter(v, 7, false, true)
}

func (mask *PPUMASK) shouldRender() bool {
  return mask.showBg || mask.showSprites
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
    result |= (0x1 << 5)
  }
  if status.Sprite0Hit {
    result |= (0x1 << 6)
  }
  if status.VBlankStarted {
    result |= (0x1 << 7)
  }

  result |= (status.LastWrite & 0x1f)

  status.VBlankStarted = false
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
    //log.Printf("Wrote PPUADDR LOW %x", v)
    //addr.TAddr |= ((uint16(v) & 0x3f) << 8)
    //addr.TAddr |= ((uint16(v) & 0x3f) << 8)
    addr.TAddr = (addr.TAddr & 0x80FF) | ((uint16(v) & 0x3F) << 8)
    //addr.TAddr &= 0x7fff
    //addr.TAddr = (addr.TAddr & 0x80FF) | ((uint16(v) & 0x3F) << 8)
    addr.WriteHi = true
  } else {
    //log.Printf("Wrote PPUADDR HI %x", v)
    //addr.TAddr |= uint16(v)
    addr.TAddr = (addr.TAddr & 0xFF00) | uint16(v)
    addr.VAddr = addr.TAddr
    addr.WriteHi = false
    //log.Printf("Wrote PPUADDR %x", addr.VAddr)
  }
}

func (addr *PPUADDR) SetOnCTRLWrite(v uint8) {
  addr.TAddr = (addr.TAddr & 0xf3ff) | uint16(v & 0x03) << 10
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
    //addr.TAddr |= uint16(v >> 3)
    //addr.FineXScroll = v & 0x3
    addr.WriteHi = true

    addr.TAddr = (addr.TAddr & 0xFFE0) | uint16(v)
    addr.FineXScroll = v & 0x07
  } else {
    //addr.TAddr |= uint16(v & 0x03) << 12
    //addr.TAddr |= uint16(v & 0xf8) << 2
    addr.TAddr = (addr.TAddr & 0x8FFF) | ((uint16(v) & 0x07) << 12)
    addr.TAddr = (addr.TAddr & 0xFC1F) | ((uint16(v) & 0xF8) << 2)
    addr.WriteHi = false
  }
}

// http://wiki.nesdev.com/w/index.php/PPU_scrolling#Y_increment
func (addr *PPUADDR) IncrementFineY() {
  v := addr.VAddr
  var y uint16

  if (v & 0x7000) != 0x7000 {
    v += 0x1000
  } else {
    //v &= ^0x7000
    v &= 0x8FFF
    y = (v & 0x03E0) >> 5
    if (y == 29) {
      y = 0
      v ^= 0x0800
    } else {
      if y == 31 {
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

  if (v & 0x001F) == 31 {
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

  // Rendering
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


  ScanLine int
  Cycle int
  Frame uint64

  paletteData   [32]byte
  nameTableData [2048]byte
  oamData       [256]byte
  front         *image.RGBA
  back          *image.RGBA

  //lowTileByte byte
  //highTileByte byte
  //nameTableByte byte
  //tileData uint64

  v uint16 // current vram address (15 bit)
  t uint16 // temporary vram address (15 bit)
  x byte   // fine x scroll (3 bit)
  w byte   // write toggle (1 bit)
  f byte   // even/odd frame flag (1 bit)

  register byte

  // NMI flags
  nmiOccurred bool
  nmiOutput   bool
  nmiPrevious bool
  nmiDelay    byte

  // background temporary variables
  nameTableByte      byte
  attributeTableByte byte
  lowTileByte        byte
  highTileByte       byte
  tileData           uint64

  // sprite temporary variables
  spriteCount      int
  spritePatterns   [8]uint32
  spritePositions  [8]byte
  spritePriorities [8]byte
  spriteIndexes    [8]byte

  // $2000 PPUCTRL
  flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
  flagIncrement       byte // 0: add 1; 1: add 32
  flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
  flagBackgroundTable byte // 0: $0000; 1: $1000
  flagSpriteSize      byte // 0: 8x8; 1: 8x16
  flagMasterSlave     byte // 0: read EXT; 1: write EXT

  // $2001 PPUMASK
  flagGrayscale          byte // 0: color; 1: grayscale
  flagShowLeftBackground byte // 0: hide; 1: show
  flagShowLeftSprites    byte // 0: hide; 1: show
  flagShowBackground     byte // 0: hide; 1: show
  flagShowSprites        byte // 0: hide; 1: show
  flagRedTint            byte // 0: normal; 1: emphasized
  flagGreenTint          byte // 0: normal; 1: emphasized
  flagBlueTint           byte // 0:

   // $2002 PPUSTATUS
  flagSpriteZeroHit  byte
  flagSpriteOverflow byte

  // $2003 OAMADDR
  oamAddress byte

  // $2007 PPUDATA
  bufferedData byte //
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
    //NametableData:    make(Memory, 0x0800),
    PaletteData:      make(Memory, 0x0020),

    PrimaryOAMBuffer:   make([]OAMSprite, 64),
    SecondaryOAMBuffer: make([]OAMSprite,  8),

    Scanline:  240,
    Dot:       340,

    Pixels: make([]byte, 4 * 256 * 240),

    ScanLine:  240,
    Cycle: 340,
  }
}

func (ppu *PPU) Reset() {
  ppu.Cycle = 340
  ppu.ScanLine = 240
  ppu.Frame = 0
  ppu.writeControl(0)
  ppu.writeMask(0)
  ppu.writeOAMAddress(0)
}

func (ppu *PPU) Run() {
  //ppu.TickScanline()
  ppu.Step()
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
  //log.Printf("Wrote PPU data at %x: %x", ppu.ADDR.VAddr, v)
  ppu.Write8(ppu.ADDR.VAddr, v)
  ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement
}

func (ppu *PPU) ReadData() uint8 {
  //log.Printf("Read PPU Data")
  val := ppu.Read8(ppu.ADDR.VAddr)
  ppu.ADDR.VAddr += ppu.CTRL.VRAMReadIncrement
  return val
}

func (ppu *PPU) OMADMA(data []uint8) {
  copy(ppu.OAMData[:], data)
}

func (ppu *PPU) Write(addr uint16, v uint8) {
  ppu.Write8(addr, v)
}

func (ppu *PPU) Write8(addr uint16, v uint8) {
  addr = addr % 0x4000
  switch {
    // Pattern tables - for now hard mapped to CHRROM
    case addr >= 0x0000 && addr < 0x2000:
      //log.Printf("Writing on CHROM at %x", addr)
      ppu.PatternTableData.Write8(addr, v)

    case addr >= 0x2000 && addr < 0x3f00:
      //ppu.NametableData.Write8(getMirroedAddr(addr), v)
      ppu.nameTableData[getMirroedAddr(addr)] = v

    case addr >= 0x3f00 && addr < 0x4000:
      ppu.PaletteData.Write8(addr % 0x20, v)

    default:
      log.Fatalf("Invalid write to PPU at %x", addr)
  }
}

func (ppu *PPU) Read(addr uint16) uint8 {
  return ppu.Read8(addr)
}

func (ppu *PPU) Read8(addr uint16) uint8 {
  addr = addr % 0x4000
  switch {
    // Pattern tables - for now hard mapped to CHRROM
    case (addr >= 0x0000 && addr < 0x2000):
      return ppu.PatternTableData.Read8(addr)

    case addr >= 0x2000 && addr < 0x3f00:
      //return ppu.NametableData.Read8(getMirroedAddr(addr))
      return ppu.nameTableData[getMirroedAddr(addr)]

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

*/

import (
  "log"
	"encoding/gob"
	"image"
	"image/color"
  "math"
)

// MINE

func addrSetter(v uint8, bitN uint8, ifNotSet uint16, ifSet uint16) uint16 {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}

func boolSetter(v uint8, bitN uint8, ifNotSet bool, ifSet bool) bool {
  if v & (0x1 << bitN) == 0 { return ifNotSet } else { return ifSet }
}

// PPUCTRL

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

//func (ppu *PPU) writeControl(value byte) {
//  //log.Printf("WROTE CONTROL %b", value)
//	ppu.flagNameTable = (value >> 0) & 3
//	ppu.flagIncrement = (value >> 2) & 1
//	ppu.flagSpriteTable = (value >> 3) & 1
//	ppu.flagBackgroundTable = (value >> 4) & 1
//	ppu.flagSpriteSize = (value >> 5) & 1
//	ppu.flagMasterSlave = (value >> 6) & 1
//	ppu.nmiOutput = (value>>7)&1 == 1
//	ppu.nmiChange()
//	// t: ....BA.. ........ = d: ......BA
//	ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
//}

	//0-1flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	//2flagIncrement       byte // 0: add 1; 1: add 32
	//3flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	//4flagBackgroundTable byte // 0: $0000; 1: $1000
	//5flagSpriteSize      byte // 0: 8x8; 1: 8x16
	//6flagMasterSlave     byte // 0: read EXT; 1: write EXT

func (ctrl *PPUCTRL) Set(v uint8) {
  //log.Printf("WROTE CONTROL %b", v)
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

//func (ppu *PPU) LowBGTileAddr() uint16 {
//  return ppu.CTRL.BgTableAddr + uint16(ppu.NameTableLatch) * 16 + ppu.ADDR.FineY()
//}

//func (ppu *PPU) HighBGTileAddr() uint16 {
//  return ppu.LowBGTileAddr() + 8
//}

// ENDOF MINE

type PPU struct {

  // MINE
  CTRL *PPUCTRL

  // ENDOF MINE

	//Memory           // memory interface
	//console *Console // reference to parent object
  CPU *CPU
  rom *Rom
  TV *TV
  Pixels []byte

	Cycle    int    // 0-340
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter

	// storage variables
	paletteData   [32]byte
	nameTableData [2048]byte
	oamData       [256]byte
	front         *image.RGBA
	back          *image.RGBA

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

	register byte

	// NMI flags
	nmiOccurred bool
	//nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	// background temporary variables
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64

	// sprite temporary variables
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	// $2000 PPUCTRL
	//flagNameTable       byte // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	//flagIncrement       byte // 0: add 1; 1: add 32
	//flagSpriteTable     byte // 0: $0000; 1: $1000; ignored in 8x16 mode
	//flagBackgroundTable byte // 0: $0000; 1: $1000
	//flagSpriteSize      byte // 0: 8x8; 1: 8x16
	//flagMasterSlave     byte // 0: read EXT; 1: write EXT

	// $2001 PPUMASK
	flagGrayscale          byte // 0: color; 1: grayscale
	flagShowLeftBackground byte // 0: hide; 1: show
	flagShowLeftSprites    byte // 0: hide; 1: show
	flagShowBackground     byte // 0: hide; 1: show
	flagShowSprites        byte // 0: hide; 1: show
	flagRedTint            byte // 0: normal; 1: emphasized
	flagGreenTint          byte // 0: normal; 1: emphasized
	flagBlueTint           byte // 0: normal; 1: emphasized

	// $2002 PPUSTATUS
	flagSpriteZeroHit  byte
	flagSpriteOverflow byte

	// $2003 OAMADDR
	oamAddress byte

	// $2007 PPUDATA
	bufferedData byte // for buffered reads
}

//func NewPPU(console *Console) *PPU {
func NewPPU(cpu *CPU, rom *Rom) *PPU {
	//ppu := PPU{Memory: NewPPUMemory(console), console: console}
	ppu := PPU{
    CPU: cpu,
    rom: rom,
    CTRL: &PPUCTRL{},
  }
	ppu.front = image.NewRGBA(image.Rect(0, 0, 256, 240))
	ppu.back = image.NewRGBA(image.Rect(0, 0, 256, 240))
  ppu.Pixels = make([]byte, 4 * 256 * 240)
	ppu.Reset()
	return &ppu
}

func (ppu *PPU) Save(encoder *gob.Encoder) error {
	encoder.Encode(ppu.Cycle)
	encoder.Encode(ppu.ScanLine)
	encoder.Encode(ppu.Frame)
	encoder.Encode(ppu.paletteData)
	encoder.Encode(ppu.nameTableData)
	encoder.Encode(ppu.oamData)
	encoder.Encode(ppu.v)
	encoder.Encode(ppu.t)
	encoder.Encode(ppu.x)
	encoder.Encode(ppu.w)
	encoder.Encode(ppu.f)
	encoder.Encode(ppu.register)
	encoder.Encode(ppu.nmiOccurred)
	//encoder.Encode(ppu.nmiOutput)
	encoder.Encode(ppu.nmiPrevious)
	encoder.Encode(ppu.nmiDelay)
	encoder.Encode(ppu.nameTableByte)
	encoder.Encode(ppu.attributeTableByte)
	encoder.Encode(ppu.lowTileByte)
	encoder.Encode(ppu.highTileByte)
	encoder.Encode(ppu.tileData)
	encoder.Encode(ppu.spriteCount)
	encoder.Encode(ppu.spritePatterns)
	encoder.Encode(ppu.spritePositions)
	encoder.Encode(ppu.spritePriorities)
	encoder.Encode(ppu.spriteIndexes)
	//encoder.Encode(ppu.flagNameTable)
	//encoder.Encode(ppu.flagIncrement)
	//encoder.Encode(ppu.flagSpriteTable)
	//encoder.Encode(ppu.flagBackgroundTable)
	//encoder.Encode(ppu.flagSpriteSize)
	//encoder.Encode(ppu.flagMasterSlave)
	encoder.Encode(ppu.flagGrayscale)
	encoder.Encode(ppu.flagShowLeftBackground)
	encoder.Encode(ppu.flagShowLeftSprites)
	encoder.Encode(ppu.flagShowBackground)
	encoder.Encode(ppu.flagShowSprites)
	encoder.Encode(ppu.flagRedTint)
	encoder.Encode(ppu.flagGreenTint)
	encoder.Encode(ppu.flagBlueTint)
	encoder.Encode(ppu.flagSpriteZeroHit)
	encoder.Encode(ppu.flagSpriteOverflow)
	encoder.Encode(ppu.oamAddress)
	encoder.Encode(ppu.bufferedData)
	return nil
}

func (ppu *PPU) Load(decoder *gob.Decoder) error {
	decoder.Decode(&ppu.Cycle)
	decoder.Decode(&ppu.ScanLine)
	decoder.Decode(&ppu.Frame)
	decoder.Decode(&ppu.paletteData)
	decoder.Decode(&ppu.nameTableData)
	decoder.Decode(&ppu.oamData)
	decoder.Decode(&ppu.v)
	decoder.Decode(&ppu.t)
	decoder.Decode(&ppu.x)
	decoder.Decode(&ppu.w)
	decoder.Decode(&ppu.f)
	decoder.Decode(&ppu.register)
	decoder.Decode(&ppu.nmiOccurred)
	//decoder.Decode(&ppu.nmiOutput)
	decoder.Decode(&ppu.nmiPrevious)
	decoder.Decode(&ppu.nmiDelay)
	decoder.Decode(&ppu.nameTableByte)
	decoder.Decode(&ppu.attributeTableByte)
	decoder.Decode(&ppu.lowTileByte)
	decoder.Decode(&ppu.highTileByte)
	decoder.Decode(&ppu.tileData)
	decoder.Decode(&ppu.spriteCount)
	decoder.Decode(&ppu.spritePatterns)
	decoder.Decode(&ppu.spritePositions)
	decoder.Decode(&ppu.spritePriorities)
	decoder.Decode(&ppu.spriteIndexes)
	//decoder.Decode(&ppu.flagNameTable)
	//decoder.Decode(&ppu.flagIncrement)
	//decoder.Decode(&ppu.flagSpriteTable)
	//decoder.Decode(&ppu.flagBackgroundTable)
	//decoder.Decode(&ppu.flagSpriteSize)
	//decoder.Decode(&ppu.flagMasterSlave)
	decoder.Decode(&ppu.flagGrayscale)
	decoder.Decode(&ppu.flagShowLeftBackground)
	decoder.Decode(&ppu.flagShowLeftSprites)
	decoder.Decode(&ppu.flagShowBackground)
	decoder.Decode(&ppu.flagShowSprites)
	decoder.Decode(&ppu.flagRedTint)
	decoder.Decode(&ppu.flagGreenTint)
	decoder.Decode(&ppu.flagBlueTint)
	decoder.Decode(&ppu.flagSpriteZeroHit)
	decoder.Decode(&ppu.flagSpriteOverflow)
	decoder.Decode(&ppu.oamAddress)
	decoder.Decode(&ppu.bufferedData)
	return nil
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
	//ppu.writeControl(0)
  ppu.t = (ppu.t & 0xF3FF) | ((uint16(0) & 0x03) << 10)
  ppu.CTRL.Set(0)
	ppu.writeMask(0)
	ppu.writeOAMAddress(0)
}

func (ppu *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return ppu.paletteData[address]
}

func (ppu *PPU) writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	ppu.paletteData[address] = value
}

func (ppu *PPU) readRegister(address uint16) byte {
  //log.Printf("Reading ppu reg %x", address)
	switch address {
	case 0x2002:
		return ppu.readStatus()
	case 0x2004:
		return ppu.readOAMData()
	case 0x2007:
		return ppu.readData()
	}
  panic("Invalid read")
	return 0
}

func (ppu *PPU) writeRegister(address uint16, value byte) {
  //log.Printf("Writing ppu reg %x: %b", address, value)
	ppu.register = value
	switch address {
	case 0x2000:
		//ppu.writeControl(value)
    ppu.CTRL.Set(value)
    ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
	case 0x2001:
		ppu.writeMask(value)
	case 0x2003:
		ppu.writeOAMAddress(value)
	case 0x2004:
		ppu.writeOAMData(value)
	case 0x2005:
		ppu.writeScroll(value)
	case 0x2006:
		ppu.writeAddress(value)
	case 0x2007:
		ppu.writeData(value)
	case 0x4014:
		ppu.writeDMA(value)
	}
}

// $2000: PPUCTRL
//func (ppu *PPU) writeControl(value byte) {
//  //log.Printf("WROTE CONTROL %b", value)
//	ppu.flagNameTable = (value >> 0) & 3
//	ppu.flagIncrement = (value >> 2) & 1
//	ppu.flagSpriteTable = (value >> 3) & 1
//	ppu.flagBackgroundTable = (value >> 4) & 1
//	ppu.flagSpriteSize = (value >> 5) & 1
//	ppu.flagMasterSlave = (value >> 6) & 1
//	ppu.nmiOutput = (value>>7)&1 == 1
//	ppu.nmiChange()
//	// t: ....BA.. ........ = d: ......BA
//	ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
//}

// $2001: PPUMASK
func (ppu *PPU) writeMask(value byte) {
  //log.Printf("WROTE MASK %b", value)
	ppu.flagGrayscale = (value >> 0) & 1
	ppu.flagShowLeftBackground = (value >> 1) & 1
	ppu.flagShowLeftSprites = (value >> 2) & 1
	ppu.flagShowBackground = (value >> 3) & 1
	ppu.flagShowSprites = (value >> 4) & 1
	ppu.flagRedTint = (value >> 5) & 1
	ppu.flagGreenTint = (value >> 6) & 1
	ppu.flagBlueTint = (value >> 7) & 1
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	result := ppu.register & 0x1F
	result |= ppu.flagSpriteOverflow << 5
	result |= ppu.flagSpriteZeroHit << 6
	if ppu.nmiOccurred {
		result |= 1 << 7
	}
	ppu.nmiOccurred = false
	ppu.nmiChange()
	// w:                   = 0
	ppu.w = 0
	return result
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddress(value byte) {
	ppu.oamAddress = value
}

// $2004: OAMDATA (read)
func (ppu *PPU) readOAMData() byte {
	return ppu.oamData[ppu.oamAddress]
}

// $2004: OAMDATA (write)
func (ppu *PPU) writeOAMData(value byte) {
	ppu.oamData[ppu.oamAddress] = value
	ppu.oamAddress++
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(value byte) {
	if ppu.w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		ppu.t = (ppu.t & 0xFFE0) | (uint16(value) >> 3)
		ppu.x = value & 0x07
		ppu.w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		ppu.w = 0
	}
}

// $2006: PPUADDR
func (ppu *PPU) writeAddress(value byte) {
	if ppu.w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		ppu.t = (ppu.t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		ppu.w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		ppu.t = (ppu.t & 0xFF00) | uint16(value)
		ppu.v = ppu.t
		ppu.w = 0
    log.Printf("Wrote PPUADDR %x", ppu.v)
	}
}

// $2007: PPUDATA (read)
func (ppu *PPU) readData() byte {
  log.Printf("Will read ppu.v = %x", ppu.v)
	value := ppu.Read(ppu.v)
	// emulate buffered reads
	if ppu.v%0x4000 < 0x3F00 {
		buffered := ppu.bufferedData
		ppu.bufferedData = value
		value = buffered
	} else {
		ppu.bufferedData = ppu.Read(ppu.v - 0x1000)
	}
	// increment address
	//if ppu.flagIncrement == 0 {
	//	ppu.v += 1
	//} else {
	//	ppu.v += 32
	//}
  ppu.v += ppu.CTRL.VRAMReadIncrement
	return value
}

// $2007: PPUDATA (write)
func (ppu *PPU) writeData(value byte) {
  //log.Printf("Wrote PPU data at %x: %x", ppu.v, value)
	ppu.Write(ppu.v, value)
	//if ppu.flagIncrement == 0 {
	//	ppu.v += 1
	//} else {
	//	ppu.v += 32
	//}
  ppu.v += ppu.CTRL.VRAMReadIncrement
}

// $4014: OAMDMA
func (ppu *PPU) writeDMA(value byte) {
	//cpu := ppu.console.CPU
	//address := uint16(value) << 8
	//for i := 0; i < 256; i++ {
	//	ppu.oamData[ppu.oamAddress] = cpu.Read(address)
	//	ppu.oamAddress++
	//	address++
	//}
	//cpu.stall += 513
	//if cpu.Cycles%2 == 1 {
	//	cpu.stall++
	//}
}

// NTSC Timing Helper Functions

func (ppu *PPU) incrementX() {
	// increment hori(v)
	// if coarse X == 31
	if ppu.v&0x001F == 31 {
		// coarse X = 0
		ppu.v &= 0xFFE0
		// switch horizontal nametable
		ppu.v ^= 0x0400
	} else {
		// increment coarse X
		ppu.v++
	}
}

func (ppu *PPU) incrementY() {
	// increment vert(v)
	// if fine Y < 7
	if ppu.v&0x7000 != 0x7000 {
		// increment fine Y
		ppu.v += 0x1000
	} else {
		// fine Y = 0
		ppu.v &= 0x8FFF
		// let y = coarse Y
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		ppu.v = (ppu.v & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU) copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

func (ppu *PPU) copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

func (ppu *PPU) nmiChange() {
	//nmi := ppu.nmiOutput && ppu.nmiOccurred
	nmi := ppu.CTRL.NMIonVBlank && ppu.nmiOccurred
	if nmi && !ppu.nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		ppu.nmiDelay = 15
	}
	ppu.nmiPrevious = nmi
}

func (ppu *PPU) setVerticalBlank() {
  //log.Printf("Will set vblank")
	ppu.front, ppu.back = ppu.back, ppu.front
  _ = math.Pow

  // COpy image to pixels
  for l := 0; l < 240; l++ {
    for c := 0; c < 256; c++ {
      r, g, b, _ := ppu.front.At(c, l).RGBA()
      pos := 4*(l*256 + c)

      // The 0 positiion seems to have no effect in any value
      ppu.Pixels[pos + 0] = 0xff//uint8(mean)
      ppu.Pixels[pos + 1] = uint8(r & 0xff)
      ppu.Pixels[pos + 2] = uint8(g & 0xff)
      ppu.Pixels[pos + 3] = uint8(b & 0xff)
    }
  }

  ppu.TV.SetFrame(ppu.Pixels)

	ppu.nmiOccurred = true
	ppu.nmiChange()
}

func (ppu *PPU) clearVerticalBlank() {
	ppu.nmiOccurred = false
	ppu.nmiChange()
}

func (ppu *PPU) fetchNameTableByte() {
	v := ppu.v
	address := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte = ppu.Read(address)
}

func (ppu *PPU) fetchAttributeTableByte() {
	v := ppu.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	ppu.attributeTableByte = ((ppu.Read(address) >> shift) & 3) << 2
}

func (ppu *PPU) fetchLowTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := 1//ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.lowTileByte = ppu.Read(address)
}

func (ppu *PPU) fetchHighTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := 1//ppu.flagBackgroundTable
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.highTileByte = ppu.Read(address + 8)
}

func (ppu *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		//a := ppu.attributeTableByte
		p1 := (ppu.lowTileByte & 0x80) >> 7
		p2 := (ppu.highTileByte & 0x80) >> 6
		ppu.lowTileByte <<= 1
		ppu.highTileByte <<= 1
		data <<= 4
		//data |= uint32(a | p1 | p2)
		data |= uint32(p1 | p2)
	}
	ppu.tileData |= uint64(data)
}

func (ppu *PPU) fetchTileData() uint32 {
	return uint32(ppu.tileData >> 32)
}

func (ppu *PPU) backgroundPixel() byte {
  _ = log.Printf
	if ppu.flagShowBackground == 0 {
		//return 0
	}
	//data := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
	data := ppu.fetchTileData()
  //log.Printf("BG Pixel %x", data)
	return byte(data & 0x0F)
}

func (ppu *PPU) spritePixel() (byte, byte) {
	if ppu.flagShowSprites == 0 {
		return 0, 0
	}
	for i := 0; i < ppu.spriteCount; i++ {
		offset := (ppu.Cycle - 1) - int(ppu.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((ppu.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1
	y := ppu.ScanLine
	//background := ppu.backgroundPixel()
  //ccc := color.RGBA{R: background, G: 0xff, A: 0xaa}
  //ppu.back.SetRGBA(x, y, ccc)
  //return
  //log.Printf("VAddr: %x", ppu.v)
	background := ppu.backgroundPixel()
	//i, sprite := ppu.spritePixel()
	//if x < 8 && ppu.flagShowLeftBackground == 0 {
	//	background = 0
	//}
	//if x < 8 && ppu.flagShowLeftSprites == 0 {
	//	sprite = 0
	//}
	//b := background%4 != 0
	//s := sprite%4 != 0
	//var color byte
	//if !b && !s {
	//	color = 0
	//} else if !b && s {
	//	color = sprite | 0x10
	//} else if b && !s {
	//	color = background
	//} else {
	//	if ppu.spriteIndexes[i] == 0 && x < 255 {
	//		ppu.flagSpriteZeroHit = 1
	//	}
	//	if ppu.spritePriorities[i] == 0 {
	//		color = sprite | 0x10
	//	} else {
	//		color = background
	//	}
	//}
  //log.Printf("BG COLOR %x", background)
	//var color byte
  //color = background
	//c := Palette[ppu.readPalette(uint16(color))%64]
	//ppu.back.SetRGBA(x, y, c)
  //log.Printf("BG PIXEL: %x", background)
  cc := color.RGBA{40*background, 40*background, 40*background, 0xff}
  ppu.back.SetRGBA(x, y, cc)
}

func (ppu *PPU) fetchSpritePattern(i, row int) uint32 {
  return 0
  /*
	tile := ppu.oamData[i*4+1]
	attributes := ppu.oamData[i*4+2]
	var address uint16
	if ppu.flagSpriteSize == 0 {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		table := ppu.flagSpriteTable
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := ppu.Read(address)
	highTileByte := ppu.Read(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
  */
}

func (ppu *PPU) evaluateSprites() {
  /*
	var h int
	if ppu.flagSpriteSize == 0 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := ppu.oamData[i*4+0]
		a := ppu.oamData[i*4+2]
		x := ppu.oamData[i*4+3]
		row := ppu.ScanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			ppu.spritePatterns[count] = ppu.fetchSpritePattern(i, row)
			ppu.spritePositions[count] = x
			ppu.spritePriorities[count] = (a >> 5) & 1
			ppu.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		ppu.flagSpriteOverflow = 1
	}
	ppu.spriteCount = count
  */
}

// tick updates Cycle, ScanLine and Frame counters
func (ppu *PPU) tick() {
	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		//if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
		if ppu.nmiDelay == 0 && ppu.CTRL.NMIonVBlank && ppu.nmiOccurred {
			//ppu.console.CPU.triggerNMI()
      ppu.CPU.nmiRequested = true
		}
	}

	//if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
	//	if ppu.f == 1 && ppu.ScanLine == 261 && ppu.Cycle == 339 {
	//		ppu.Cycle = 0
	//		ppu.ScanLine = 0
	//		ppu.Frame++
	//		ppu.f ^= 1
	//		return
	//	}
	//}
	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.ScanLine++
		if ppu.ScanLine > 261 {
			ppu.ScanLine = 0
			//ppu.Frame++
			//ppu.f ^= 1
		}
	}
}

// Step executes a single PPU cycle
func (ppu *PPU) Step() {
	ppu.tick()

	renderingEnabled := ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0
	preLine := ppu.ScanLine == 261
	visibleLine := ppu.ScanLine < 240
	// postLine := ppu.ScanLine == 240
	renderLine := preLine || visibleLine
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// background logic
	if renderingEnabled {
		if visibleLine && visibleCycle {
			ppu.renderPixel()
		}
		if renderLine && fetchCycle {
      //log.Printf("tileData: %x", ppu.tileData)
			ppu.tileData <<= 4
			switch ppu.Cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				//ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchLowTileByte()
			case 7:
				ppu.fetchHighTileByte()
			case 0:
				ppu.storeTileData()
			}
		}
		if preLine && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
			ppu.copyY()
		}
		if renderLine {
			if fetchCycle && ppu.Cycle%8 == 0 {
				ppu.incrementX()
			}
			if ppu.Cycle == 256 {
				ppu.incrementY()
			}
			if ppu.Cycle == 257 {
				ppu.copyX()
			}
		}
	}

	// sprite logic
	//if renderingEnabled {
	//	if ppu.Cycle == 257 {
	//		if visibleLine {
	//			ppu.evaluateSprites()
	//		} else {
	//			ppu.spriteCount = 0
	//		}
	//	}
	//}

	// vblank logic
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.setVerticalBlank()
    //ppu.console.CPU.triggerNMI()
	}
	if preLine && ppu.Cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.flagSpriteZeroHit = 0
		ppu.flagSpriteOverflow = 0
	}
}

// From memory

//func (mem *ppuMemory) Read(address uint16) byte {
func (ppu *PPU) Read(address uint16) byte {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		//return mem.console.Mapper.Read(address)
    return ppu.rom.CHRROM[address]
	case address < 0x3F00:
		//mode := mem.console.Cartridge.Mirror
		//return mem.console.PPU.nameTableData[MirrorAddress(mode, address)%2048]
    return ppu.nameTableData[address%0x800]
	case address < 0x4000:
		return ppu.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
	}
	return 0
}

//func (mem *ppuMemory) Write(address uint16, value byte) {
func (ppu *PPU) Write(address uint16, value byte) {
	address = address % 0x4000
	switch {
	case address < 0x2000:
		//mem.console.Mapper.Write(address, value)
    ppu.rom.CHRROM[address] = value
	case address < 0x3F00:
		//mode := mem.console.Cartridge.Mirror
		//mem.console.PPU.nameTableData[MirrorAddress(mode, address)%2048] = value
    ppu.nameTableData[address%0x800] = value
	case address < 0x4000:
		ppu.writePalette(address%32, value)
	default:
		log.Fatalf("unhandled ppu memory write at address: 0x%04X", address)
	}
}

//func (ppu *PPU) Write8(addr uint16, v uint8) {
//  addr = addr % 0x4000
//  switch {
//    // Pattern tables - for now hard mapped to CHRROM
//    case addr >= 0x0000 && addr < 0x2000:
//      //log.Printf("Writing on CHROM at %x", addr)
//      ppu.PatternTableData.Write8(addr, v)
//
//    case addr >= 0x2000 && addr < 0x3f00:
//      //ppu.NametableData.Write8(getMirroedAddr(addr), v)
//      ppu.nameTableData[getMirroedAddr(addr)] = v
//
//    case addr >= 0x3f00 && addr < 0x4000:
//      ppu.PaletteData.Write8(addr % 0x20, v)
//
//    default:
//      log.Fatalf("Invalid write to PPU at %x", addr)
//  }
//}

// Mirroring Modes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
	MirrorFour       = 4
)

var MirrorLookup = [...][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 1, 2, 3},
}

func MirrorAddress(mode byte, address uint16) uint16 {
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}
