package awesomenes

//import (
//  "log"
//  "math/rand"
//)

/*
  Screen resolution: 256 cols * 240 rows pixels
  Scanlines: 262 per frame
  Dots:      341 per scanline

  Timings extracted from http://wiki.nesdev.com/w/images/d/d1/Ntsc_timing.png
*/

/*
func makeRandom(pixels []byte) {
  for i := 0; i < 240; i++ {
    for j := 0; j < 256; j++ {
      pixels[i*240 + j] = uint8(rand.Uint32() % 32)
    }
  }
}

func (ppu *PPU) TickScanline() {
  line := ppu.Scanline
  //log.Printf("Scanline %v", line)
  lineType := scanlineType(line)

  // Pre-render scanline
  if lineType == SCANLINE_TYPE_PRE {
    ppu.tickPreScanline()

  // Visible scanline
  } else if lineType == SCANLINE_TYPE_VISIBLE {
    ppu.tickVisibleScanline()

  } else if line == SCANLINE_NMI {
    if ppu.Dot == 1 {
      //log.Printf("VBLANK WILL START\n")
      ppu.STATUS.VBlankStarted = true
      //makeRandom(ppu.Pixels)
      //ppu.tv.SetFrame(ppu.Pixels)
      if ppu.CTRL.NMIonVBlank {
        ppu.CPU.nmiRequested = true
      }
    }
  } else if lineType == SCANLINE_TYPE_POST {
    if ppu.Dot == 0 {
      ppu.tv.SetFrame(ppu.Pixels)
    }
  }

  //log.Printf("Line: %v", line)
  ppu.Dot += 1
  if ppu.Dot == 341 {
    ppu.Scanline += 1
    if ppu.Scanline == 262 {
      // Wrap around
      ppu.Scanline = 0
    }
    ppu.Dot = 0
  }
}

func (ppu *PPU) tickPreScanline() {
  dot := ppu.Dot

  if dot == 1 {
    //Not in VBlank anymore. Prepare for next visible scanlines.
    ppu.STATUS.VBlankStarted  = false
    ppu.STATUS.Sprite0Hit     = false
    ppu.STATUS.SpriteOverflow = false

  } else if dot >= 280 && dot <= 304 {
    if ppu.MASK.shouldRender() {
      ppu.ADDR.TransferY()
    }
  }

  // Now do everything a visible line does
  ppu.tickVisibleScanline()
}

func (ppu *PPU) tickVisibleScanline() {
  dot := ppu.Dot

  // Background evaluation
  if !ppu.MASK.shouldRender() {
    return
  }

  if dot >= 1 && dot <= 255 {
    ppu.RenderSinglePixel()
  }

  //log.Printf("ppu.tileData: %x", ppu.tileData)
  if (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 340) {
     //ppu.tileData <<= 4
     switch dot % 8 {
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
    //switch ppu.Dot % 8 {
    //  case 1:
    //    ppu.BgTileShiftLow  |= uint16(ppu.BgLatchLow)
    //    ppu.BgTileShiftHigh |= uint16(ppu.BgLatchHigh)
    //    ppu.storeTileData()
    //  case 2:
    //    ppu.NameTableLatch = ppu.Read8(ppu.ADDR.NameTableAddr())
    //  case 4:
    //    ppu.AttrTableLatch = ppu.Read8(ppu.ADDR.AttrTableAddr())
    //  case 5:
    //    //log.Printf("Low bg tile addr %x", ppu.LowBGTileAddr())
    //    ppu.BgLatchLow = ppu.Read8(ppu.LowBGTileAddr())
    //    //log.Printf("LathLow %x", ppu.BgLatchLow)
    //  case 7:
    //    //log.Printf("Low bg tile addr %x", ppu.HighBGTileAddr())
    //    ppu.BgLatchHigh = ppu.Read8(ppu.HighBGTileAddr())
    //}
  }

  // Sprite evaluation

  //if dot == 1 {
  //  ppu.ClearSecondaryOAM()
  //} else if dot == 256 {
  //  ppu.EvalSprites()
  //}

  //// Housekeeping. See  http://wiki.nesdev.com/w/index.php/PPU_scrolling

  if dot == 256 {
    ppu.ADDR.IncrementFineY()
  }

  if dot == 257 {
    ppu.ADDR.TransferX()
  }

  //fetchCycle := (dot >= 321 && dot <= 336) || (dot >= 1 && dot <= 256)
  //if fetchCycle && dot % 8 == 0 {
  if dot % 8 == 0 {
    ppu.ADDR.IncrementCoarseX()
  }

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
  //data := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
  data := ppu.fetchTileData()
  //log.Printf("BG Pixel %x", data)
  return byte(data & 0x0F)
}

func (ppu *PPU) fetchNameTableByte() {
  v := ppu.ADDR.VAddr
  address := 0x2000 | (v & 0x0FFF)
  ppu.nameTableByte = ppu.Read8(address)
}


func (ppu *PPU) fetchLowTileByte() {
  fineY := (ppu.ADDR.VAddr >> 12) & 7
  table := 1//ppu.flagBackgroundTable
  tile := ppu.nameTableByte
  address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
  ppu.lowTileByte = ppu.Read8(address)
  //log.Printf("GOT TILE LOW: %x", ppu.lowTileByte)
}

func (ppu *PPU) fetchHighTileByte() {
  //log.Printf("GOT TILE HIGH : %x", ppu.highTileByte)
  fineY := (ppu.ADDR.VAddr >> 12) & 7
  table := 1//ppu.flagBackgroundTable
  tile := ppu.nameTableByte
  address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
  ppu.highTileByte = ppu.Read8(address + 8)
}



func (ppu *PPU) RenderSinglePixel() {
  line := ppu.Scanline
  dot := ppu.Dot

  //log.Printf("Rendering line %v dot %v", line, dot)
  //log.Printf("VAddr: %x", ppu.ADDR.VAddr)

  //test := ((ppu.BgTileShiftHigh >> 15) << 1) | (ppu.BgTileShiftLow >> 15)
  //v := uint8((80 * test) & 0xff)

  //if line < 50 {
  //  v = 0xff
  //} else {
  //  v = 0x00
  //}
  //v = uint8(80 * (test & 0xff))
  v := 40*ppu.backgroundPixel()


  //log.Printf("BG PIXEL %x", v)

  //ppu.Pixels[line * 255 + dot] = 60*uint8(test)
  //if line >= 0 && line <= 239 {
    //v := 80 * uint8(rand.Uint32() & 0x03)
    //vv := v << 6 | v << 4 | v << 2 | 0x3
    ppu.Pixels[3*(line * 256 + dot) + 0] = v
    ppu.Pixels[3*(line * 256 + dot) + 1] = v
    ppu.Pixels[3*(line * 256 + dot) + 2] = v
    //ppu.Pixels[4*(line * 255 + dot) + 3] = 0xff
  //}

  //if line >= 0 && line <= 239 {
  //  //ppu.Pixels[(line * 256 + dot + 0)] = 40*uint8(test)
  //  //ppu.Pixels[(line * 256 + dot + 1)] = 40*uint8(test)
  //  //ppu.Pixels[(line * 256 + dot + 2)] = 40*uint8(test)
  //  //ppu.Pixels[(line * 256 + dot + 3)] = 0xff
  //  //ppu.Pixels[1*(line * 255 + dot) + 0] = uint8(line)
  //  //ppu.Pixels[1*(line * 255 + dot) + 1] = uint8(line)
  //  //ppu.Pixels[1*(line * 255 + dot) + 2] = uint8(line)
  //  //ppu.Pixels[1*(line * 255 + dot) + 3] = uint8(line)
  //  //ppu.Pixels[(line * 256 + dot + 1)] = 40*uint8(test & 0x1)
  //  //ppu.Pixels[(line * 256 + dot + 2)] = 40*uint8(test & 0x1)
  //  //ppu.Pixels[(line * 256 + dot + 3)] = 40*uint8(test & 0x1)
  //}

  //log.Printf("Test: %x line %v", test, line)

  ppu.BgTileShiftHigh <<= 1
  ppu.BgTileShiftLow  <<= 1
}

// Noop is fine?
func (ppu *PPU) ClearSecondaryOAM() {
  return
}

func (ppu *PPU) EvalSprites() {
  return
}

const (
  SCANLINE_TYPE_PRE     = 0x1
  SCANLINE_TYPE_VISIBLE = 0x2
  SCANLINE_TYPE_POST    = 0x3
  SCANLINE_TYPE_VBLANK  = 0x4

  SCANLINE_NMI          = 241

  DOT_TYPE_VISIBLE      = 0x1
  DOT_TYPE_PREFETCH     = 0x2
  DOT_TYPE_INVISIBLE    = 0x3
)

func scanlineType(scanlineN int) int {
  switch {
    case scanlineN == 261:
      return SCANLINE_TYPE_PRE

    case scanlineN < 240:
      return SCANLINE_TYPE_VISIBLE

    case scanlineN == 240:
      return SCANLINE_TYPE_POST

    case scanlineN >= 241 && scanlineN <= 260:
      return SCANLINE_TYPE_VBLANK

    default:
      log.Fatalf("Invalid scanline number %d\n", scanlineN)
      return 0
  }
}

func DotType(dot int) int {
  switch {
    case dot > 1 && dot <= 256:
      return DOT_TYPE_VISIBLE

    case dot >= 257 && dot <= 336:
      return DOT_TYPE_PREFETCH

    default:
      return DOT_TYPE_INVISIBLE
  }
}


func (ppu *PPU) Step2() {
  ppu.tick()

  renderingEnabled := ppu.MASK.shouldRender()
  preLine := ppu.Scanline == 261
  visibleLine := ppu.Scanline < 240
  // postLine := ppu.ScanLine == 240
  renderLine := preLine || visibleLine
  preFetchDot := ppu.Dot >= 321 && ppu.Dot <= 336
  visibleDot := ppu.Dot >= 1 && ppu.Dot <= 256
  fetchDot := preFetchDot || visibleDot

  // background logic
  if renderingEnabled {
    if visibleLine && visibleDot {
      //ppu.renderPixel()
      ppu.RenderSinglePixel()
    }
    if renderLine && fetchDot {
      //log.Printf("tileData: %x", ppu.tileData)
      ppu.tileData <<= 4
      switch ppu.Dot % 8 {
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
    if preLine && ppu.Dot >= 280 && ppu.Dot <= 304 {
      ppu.ADDR.TransferY()
    }
    if renderLine {
      if fetchDot && ppu.Dot%8 == 0 {
        ppu.ADDR.IncrementCoarseX()
      }
      if ppu.Dot == 256 {
        ppu.ADDR.IncrementFineY()
      }
      if ppu.Dot == 257 {
        ppu.ADDR.TransferX()
      }
    }
  }

  // vblank logic
  if ppu.Scanline == 241 && ppu.Dot == 1 {
    ppu.STATUS.VBlankStarted = true
    if ppu.CTRL.NMIonVBlank {
      ppu.CPU.nmiRequested = true
    }
  }
  if preLine && ppu.Dot == 1 {
    ppu.STATUS.VBlankStarted  = false
    ppu.STATUS.Sprite0Hit     = false
    ppu.STATUS.SpriteOverflow = false
  }
}

func (ppu *PPU) tick2() {
  //if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
  //  if ppu.f == 1 && ppu.ScanLine == 261 && ppu.Dot == 339 {
  //    ppu.Dot = 0
  //    ppu.ScanLine = 0
  //    ppu.Frame++
  //    ppu.f ^= 1
  //    return
  //  }
  //}
  ppu.Dot++
  if ppu.Dot > 340 {
    ppu.Dot = 0
    ppu.Scanline++
    if ppu.Scanline > 261 {
      ppu.Scanline = 0
      //ppu.Frame++
      //ppu.f ^= 1
    }
  }
}

func (ppu *PPU) writeControl(value byte) {
  //log.Printf("WROTE CONTROL %b", value)
  ppu.flagNameTable = (value >> 0) & 3
  ppu.flagIncrement = (value >> 2) & 1
  ppu.flagSpriteTable = (value >> 3) & 1
  ppu.flagBackgroundTable = (value >> 4) & 1
  ppu.flagSpriteSize = (value >> 5) & 1
  ppu.flagMasterSlave = (value >> 6) & 1
  ppu.nmiOutput = (value>>7)&1 == 1
  ppu.nmiChange()
  // t: ....BA.. ........ = d: ......BA
  ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

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
    //log.Printf("Wrote PPUADDR %x", ppu.v)
  }
}

func (ppu *PPU) readData() byte {
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
  if ppu.flagIncrement == 0 {
    ppu.v += 1
  } else {
    ppu.v += 32
  }
  return value
}

// $2007: PPUDATA (write)
func (ppu *PPU) writeData(value byte) {
  //log.Printf("Wrote PPU data at %x: %x", ppu.v, value)
  ppu.Write(ppu.v, value)
  if ppu.flagIncrement == 0 {
    ppu.v += 1
  } else {
    ppu.v += 32
  }
}

// $4014: OAMDMA
func (ppu *PPU) writeDMA(value byte) {
  cpu := ppu.CPU
  address := uint16(value) << 8
  for i := 0; i < 256; i++ {
    ppu.oamData[ppu.oamAddress] = cpu.mem.Read8(address)
    ppu.oamAddress++
    address++
  }
  //cpu.stall += 513
  //if cpu.Cycles%2 == 1 {
  //  cpu.stall++
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
  nmi := ppu.nmiOutput && ppu.nmiOccurred
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
  ppu.nmiOccurred = true
  ppu.nmiChange()
}

func (ppu *PPU) renderPixel() {
  return
  x := ppu.Cycle - 1
  y := ppu.ScanLine
  background := ppu.backgroundPixel()
  log.Printf("BG PIXEL: %x", background)
  //cc := color.RGBA{40*background, 40*background, 40*background, 0xff}
  //ppu.back.SetRGBA(x, y, cc)
  line := y
  dot := x
  v := 40*uint8(background)
  ppu.Pixels[3*(line * 256 + dot) + 0] = v
  ppu.Pixels[3*(line * 256 + dot) + 1] = v
  ppu.Pixels[3*(line * 256 + dot) + 2] = v
}

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
  //  if ppu.Cycle == 257 {
  //    if visibleLine {
  //      ppu.evaluateSprites()
  //    } else {
  //      ppu.spriteCount = 0
  //    }
  //  }
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

func (ppu *PPU) tick() {
  if ppu.nmiDelay > 0 {
    ppu.nmiDelay--
    if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
      //ppu.console.CPU.triggerNMI()
      log.Printf("WILL TRIGGER NMI")
      ppu.CPU.nmiRequested = true
    }
  }

  //if ppu.flagShowBackground != 0 || ppu.flagShowSprites != 0 {
  //  if ppu.f == 1 && ppu.ScanLine == 261 && ppu.Cycle == 339 {
  //    ppu.Cycle = 0
  //    ppu.ScanLine = 0
  //    ppu.Frame++
  //    ppu.f ^= 1
  //    return
  //  }
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

func (ppu *PPU) clearVerticalBlank() {
  ppu.nmiOccurred = false
  ppu.nmiChange()
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
    ppu.writeControl(value)
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

*/
