package awesomenes

import (
  "log"
)

/*
  Screen resolution: 256 cols * 240 rows pixels
  Scanlines: 262 per frame
  Dots:      341 per scanline

  Timings extracted from http://wiki.nesdev.com/w/images/d/d1/Ntsc_timing.png
*/

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
      ppu.setVerticalBlank()
      //ppu.STATUS.VBlankStarted = true
      //if ppu.CTRL.NMIonVBlank {
      //  ppu.CPU.nmiRequested = true
      //}
    }
  } else if lineType == SCANLINE_TYPE_POST {
    if ppu.Dot == 0 {
      //ppu.TV.SetFrame(ppu.Pixels)
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
  dot         := ppu.Dot
  isFetchTime := (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 340)

  if !ppu.MASK.shouldRender() {
    return
  }

  if dot >= 1 && dot <= 255 {
    //ppu.RenderSinglePixel()
    ppu.renderPixel()
  }

  //log.Printf("ppu.tileData: %x", ppu.tileData)
  //if (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 340) {
  if isFetchTime {
     ppu.tileData <<= 4
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

  // Housekeeping. See http://wiki.nesdev.com/w/index.php/PPU_scrolling

  if dot == 256 {
    ppu.ADDR.IncrementFineY()
  }

  if dot == 257 {
    ppu.ADDR.TransferX()
  }

  if isFetchTime && dot % 8 == 0 {
    ppu.ADDR.IncrementCoarseX()
  }

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
