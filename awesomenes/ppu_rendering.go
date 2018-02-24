package awesomenes

import "log"

// 256 * 240 pixels
// 262 scanlines rendered per frame

func (ppu *PPU) TickScanline() {
  lineType := scanlineType(ppu.Scanline)

  if lineType == SCANLINE_TYPE_PRE || lineType == SCANLINE_TYPE_VISIBLE {
    ppu.tickVisibleScanline()
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

func (ppu *PPU) tickVisibleScanline() {
  line := ppu.Scanline
  lineType := scanlineType(line)
  dot := ppu.Dot
  dotType := DotType(dot)

  if dotType == DOT_TYPE_VISIBLE {
    ppu.RenderSinglePixel()
  }

  if dotType == DOT_TYPE_VISIBLE || dotType == DOT_TYPE_PREFETCH {
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

  // http://wiki.nesdev.com/w/index.php/PPU_scrolling
  if line == 256 {
    ppu.ADDR.IncrementFineY()
  }

  if line == 257 {
    ppu.ADDR.TransferX()
  }

  if lineType == SCANLINE_TYPE_PRE && dot >= 280 && dot <= 304 {
    ppu.ADDR.TransferY()
  }

  if dotType == DOT_TYPE_VISIBLE || dotType == DOT_TYPE_PREFETCH {
  }
}

func (ppu *PPU) RenderSinglePixel() {
}

const (
  SCANLINE_TYPE_PRE     = 0x1
  SCANLINE_TYPE_VISIBLE = 0x2
  SCANLINE_TYPE_POST    = 0x3
  SCANLINE_TYPE_VBLANK  = 0x4

  SCANLINE_NMI          = 241

  DOT_TYPE_VISIBLE    = 0x1
  DOT_TYPE_PREFETCH   = 0x2
  DOT_TYPE_INVISIBLE  = 0x3
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
