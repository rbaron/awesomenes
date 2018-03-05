package awesomenes

import (
  "log"
  "image/color"
)

// From http://www.thealmightyguru.com/Games/Hacking/Wiki/index.php?title=NES_Palette
var Palette = [64]uint32 {
0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084, 0xA80020, 0xA81000, 0x881400,
0x503000, 0x007800, 0x006800, 0x005800, 0x004058, 0x000000, 0x000000, 0x000000,
0xBCBCBC, 0x0078F8, 0x0058F8, 0x6844FC, 0xD800CC, 0xE40058, 0xF83800, 0xE45C10,
0xAC7C00, 0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000, 0x000000, 0x000000,
0xF8F8F8, 0x3CBCFC, 0x6888FC, 0x9878F8, 0xF878F8, 0xF85898, 0xF87858, 0xFCA044,
0xF8B800, 0xB8F818, 0x58D854, 0x58F898, 0x00E8D8, 0x787878, 0x000000, 0x000000,
0xFCFCFC, 0xA4E4FC, 0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0, 0xFCE0A8,
0xF8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8, 0x00FCFC, 0xF8D8F8, 0x000000, 0x000000,
}

/*
  Screen resolution: 256 cols * 240 rows pixels
  Scanlines: 262 per frame
  Dots:      341 per scanline

  Timings extracted from http://wiki.nesdev.com/w/images/d/d1/Ntsc_timing.png
*/

func (ppu *PPU) TickScanline() {
  line := ppu.Scanline
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
  isFetchTime := (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 336)

  if !ppu.MASK.shouldRender() {
    return
  }

  if dot >= 1 && dot <= 256 {
    ppu.RenderSinglePixel()
  }

  if isFetchTime {

    ppu.BgTileShiftLow  <<= 1
    ppu.BgTileShiftHigh <<= 1
    ppu.AttrShiftLow    <<= 1
    ppu.AttrShiftHigh   <<= 1
    ppu.AttrShiftLow    |= (ppu.AttrLatchLow  << 0)
    ppu.AttrShiftHigh   |= (ppu.AttrLatchHigh << 1)

    switch ppu.Dot % 8 {
      case 1:
        ppu.tempTileAddr    = ppu.ADDR.NameTableAddr()

        // Feed new data into the background tile latches
        ppu.BgTileShiftLow  |= uint16(ppu.BgLatchLow)
        ppu.BgTileShiftHigh |= uint16(ppu.BgLatchHigh)

        // Feed new data into the attribute latches
        ppu.AttrLatchLow    = (ppu.AttrTableLatch >> 0) & 0x1
        ppu.AttrLatchHigh   = (ppu.AttrTableLatch >> 1) & 0x1
      case 2:
        ppu.NameTableLatch  = ppu.Read(ppu.tempTileAddr)
      case 3:
        ppu.tempTileAddr    = ppu.ADDR.AttrTableAddr()
      case 4:
        shift := ((ppu.ADDR.VAddr >> 4) & 4) | (ppu.ADDR.VAddr & 2)
        ppu.AttrTableLatch  = ppu.Read(ppu.tempTileAddr) >> shift
      case 5:
        ppu.tempTileAddr    = ppu.LowBGTileAddr()
      case 6:
        ppu.BgLatchLow      = ppu.Read(ppu.tempTileAddr)
      case 7:
        ppu.tempTileAddr    = ppu.HighBGTileAddr()
      case 0:
        ppu.BgLatchHigh     = ppu.Read(ppu.tempTileAddr)
        ppu.ADDR.IncrementCoarseX()
    }
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

  //if isFetchTime && dot % 8 == 0 {
  if isFetchTime && dot % 8 == 0 {
    //ppu.ADDR.IncrementCoarseX()
  }

}

func (ppu *PPU) LowBGTileAddr() uint16 {
  return ppu.CTRL.BgTableAddr + uint16(ppu.NameTableLatch) * 16 + ppu.ADDR.FineY()
}

func (ppu *PPU) HighBGTileAddr() uint16 {
  return ppu.LowBGTileAddr() + 8
}

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

func (ppu *PPU) RenderSinglePixel() {
  x := ppu.Dot - 1
	y := ppu.Scanline


  // TODO: account for fine X scrolling
  background := uint8(
    uint16((ppu.AttrShiftHigh >> 7) << 3) |
    uint16((ppu.AttrShiftLow  >> 7) << 2) |
    ((ppu.BgTileShiftHigh >> 15) << 1)    |
    ((ppu.BgTileShiftLow  >> 15) << 0))

  if background & 0x03 == 0x0 {
    background = 0
  }

  addr := ppu.Read(0x3f00 + uint16(background))
  c := Palette[addr]

  r := uint8((c >> 16) & 0xff)
  g := uint8((c >>  8) & 0xff)
  b := uint8((c >>  0) & 0xff)

  cc := color.RGBA{r, g, b, 0x00}

  ppu.back.SetRGBA(x, y, cc)
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
