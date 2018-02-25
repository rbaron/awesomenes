package awesomenes

import "log"

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
    //log.Printf("VBLANK WILL START\n")
    ppu.STATUS.VBlankStarted = true
    if ppu.CTRL.NMIonVBlank {
      ppu.CPU.nmiRequested = true
    }
  }

  ppu.Dot += 1
  if ppu.Dot == 341 {
    ppu.Scanline += 1
    if ppu.Scanline == 262 {
      // Frame is ready for displaying

      // Trigger nmi?
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
    ppu.ADDR.TransferY()
  }

  // Now do everything a visible line does
  ppu.tickVisibleScanline()
}

func (ppu *PPU) tickVisibleScanline() {
  dot := ppu.Dot

  // Background evaluation

  if (dot >= 1 && dot <= 256) || (dot >= 321 && dot <= 340) {
    switch ppu.Dot % 8 {
      case 1:
        ppu.BgTileShiftLow  |= uint16(ppu.BgLatchLow)
        ppu.BgTileShiftHigh |= uint16(ppu.BgLatchHigh)
      case 2:
        ppu.NameTableLatch = ppu.Read8(ppu.ADDR.NameTableAddr())
      case 4:
        ppu.AttrTableLatch = ppu.Read8(ppu.ADDR.AttrTableAddr())
      case 5:
        ppu.BgLatchLow = ppu.Read8(ppu.LowBGTileAddr())
      case 7:
        ppu.BgLatchHigh = ppu.Read8(ppu.HighBGTileAddr())
    }
  }

  // Sprite evaluation

  if dot == 1 {
    ppu.ClearSecondaryOAM()
  } else if dot == 256 {
    ppu.EvalSprites()
  }

  // Housekeeping. See  http://wiki.nesdev.com/w/index.php/PPU_scrolling

  if dot == 256 {
    ppu.ADDR.IncrementFineY()
  }

  if dot == 257 {
    ppu.ADDR.TransferX()
  }

  if dot >= 1 && dot % 8 == 0 {
    ppu.ADDR.IncrementCoarseX()
  }

  ppu.RenderSinglePixel()
}

func (ppu *PPU) RenderSinglePixel() {
  line := ppu.Scanline
  dot := ppu.Dot

  ppu.Pixels[(line * 256 + dot) % (256 * 240)] = uint8((line + dot) & 0xff)
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
