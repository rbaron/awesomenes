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
  } else if dot == 257 {
    ppu.spriteCount = 0
    ppu.spriteCount2 = 0
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


  // Background evaluation

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

        // Feed new data into the background tile shift registers
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

  if dot == 1 {
    //ppu.ClearSecondaryOAM()
  } else if dot == 257 {
    //ppu.EvalSprites()
    //ppu.EvalSprites2()
    ppu.EvalSprites2()
  }

  if dot == 256 {
    ppu.ADDR.IncrementFineY()
  }

  if dot == 257 {
    ppu.ADDR.TransferX()
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
  x  := ppu.Dot - 2
	y  := ppu.Scanline
  fx := ppu.ADDR.FineXScroll

  var background uint8

  if ppu.MASK.showBg {
    // Pull bg attribute and tile from the 4 shift registers into a 4-bit word
    background = uint8(
      uint16(((ppu.AttrShiftHigh >> (7 - fx)) & 0x1) << 3) |
      uint16(((ppu.AttrShiftLow  >> (7 - fx)) & 0x1) << 2) |
      (((ppu.BgTileShiftHigh >> (15 - fx)) & 0x1) << 1)    |
      (((ppu.BgTileShiftLow  >> (15 - fx)) & 0x1) << 0))
  }

  //i, sprite := ppu.spritePixel()
  i, sprite := ppu.spritePixel2()
  //i, sprite := ppu.spritePixel3()

  entry := ppu.oamEntries[i]
  //log.Printf("2. Got sprite id %v pixel %x\n\n", entry.id, sprite)

  //log.Printf("Got sprite pixel %d, %x", i, sprite)
  //log.Printf("entry id: %x", entry.id)

	b1 := background%4 != 0
	s := sprite%4 != 0

	var color1 byte
	if !b1 && !s {
		color1 = 0
	} else if !b1 && s {
		color1 = sprite | 0x10
	} else if b1 && !s {
		color1 = background
	} else {
	  //if ppu.spriteIndexes[i] == 0 && x < 255 {
		if entry.id == 0 && x < 255 {
      log.Printf("Will set hit")
      ppu.STATUS.Sprite0Hit = true
    }
		//if ppu.spritePriorities[i] == 0 {
		if entry.priority == 0 {
			color1 = sprite | 0x10
		} else {
			color1 = background
		}
	}

  color1 = sprite | 0x10
  addr := ppu.Read(0x3f00 + uint16(color1))
  c := Palette[addr]

  r := uint8((c >> 16) & 0xff)
  g := uint8((c >>  8) & 0xff)
  b := uint8((c >>  0) & 0xff)

  cc := color.RGBA{r, g, b, 0x00}

  ppu.back.SetRGBA(x, y, cc)
}

func unpackOAMEntry(entryIdx int, oamData []uint8) (y, tileN, attrs, x uint8) {
  baseIdx := 4 * entryIdx

  y     = oamData[baseIdx + 0]
  tileN = oamData[baseIdx + 1]
  attrs = oamData[baseIdx + 2]
  x     = oamData[baseIdx + 3]

  return
}

func (ppu *PPU) EvalSprites2() {
  nSpritesInScanline := 0

  for oamIdx := 0; oamIdx < len(ppu.oamData)/4; oamIdx++ {
    y, tileN, attrs, x := unpackOAMEntry(oamIdx, ppu.oamData[:])

    // Is this sprite in the current scanline?
		row := ppu.Scanline - int(y)
		if row < 0 || row >= int(ppu.CTRL.SpriteSizeU) {
			continue
		} else {
    //if int(y) >= ppu.Scanline && int(y) < ppu.Scanline + int(ppu.CTRL.SpriteSizeU) {
      if nSpritesInScanline < 8 {

        spriteRow := ppu.Scanline - int(y)
        addr := ppu.CTRL.SpritePatTableAddr + uint16(tileN)*16 + uint16(spriteRow)

        ppu.oamEntries[nSpritesInScanline].id       = oamIdx
        ppu.oamEntries[nSpritesInScanline].x        = x
        ppu.oamEntries[nSpritesInScanline].attrs    = attrs
        ppu.oamEntries[nSpritesInScanline].tileLow  = ppu.Read(addr)
        ppu.oamEntries[nSpritesInScanline].tileHigh = ppu.Read(addr + 8)
        //ppu.oamEntries[nSpritesInScanline].patt     = ppu.fetchSpritePattern(oamIdx, spriteRow)
        ppu.oamEntries[nSpritesInScanline].priority = (attrs >> 5) & 0x1

        nSpritesInScanline++
      } else {
        //log.Printf("OVER")
        ppu.STATUS.SpriteOverflow = true
      }
    }
  }
  ppu.spriteCount2 = nSpritesInScanline
}

func (ppu *PPU) EvalSprites3() {
	var h = int(ppu.CTRL.SpriteSizeU)

	count := 0
	for i := 0; i < 64; i++ {
		y := ppu.oamData[i*4+0]
		tile := ppu.oamData[i*4+1]
		a := ppu.oamData[i*4+2]
		x := ppu.oamData[i*4+3]
		row := ppu.Scanline - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			//ppu.spritePatterns[count] = ppu.fetchSpritePattern(i, row)
			//ppu.spritePositions[count] = x
			//ppu.spritePriorities[count] = (a >> 5) & 1
			//ppu.spriteIndexes[count] = byte(i)
      //log.Printf("Got tiles for id %v", byte(i))
        ppu.oamEntries[count].id       = i
        ppu.oamEntries[count].x        = x
        ppu.oamEntries[count].attrs    = a
        //ppu.oamEntries[count].patt     = ppu.fetchSpritePattern(i, row)
        ppu.oamEntries[count].priority = (a >> 5) & 0x1

        address := ppu.CTRL.SpritePatTableAddr + uint16(tile)*16 + uint16(row)

        ppu.oamEntries[count].tileLow  = ppu.Read(address)
        ppu.oamEntries[count].tileHigh = ppu.Read(address + 8)

        //log.Printf("Got tiles %x %x", ppu.Read(address), ppu.Read(address+ 8))
		}
		count++
	}
	if count > 8 {
		count = 8
    ppu.STATUS.SpriteOverflow = true
	}
	ppu.spriteCount2 = count
}


func (ppu *PPU) EvalSprites() {
	var h int
  if ppu.CTRL.SpriteSize == SPRITE_SIZE_8 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := ppu.oamData[i*4+0]
		a := ppu.oamData[i*4+2]
		x := ppu.oamData[i*4+3]
		row := ppu.Scanline - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			ppu.spritePatterns[count] = ppu.fetchSpritePattern(i, row)
			ppu.spritePositions[count] = x
			ppu.spritePriorities[count] = (a >> 5) & 1
			ppu.spriteIndexes[count] = byte(i)
      //log.Printf("Got tiles for id %v", byte(i))

		}
		count++
	}
	if count > 8 {
		count = 8
    ppu.STATUS.SpriteOverflow = true
	}
	ppu.spriteCount = count
}

func (ppu *PPU) fetchSpritePattern(i int, row int) uint32 {
	tile := ppu.oamData[i*4+1]
	attributes := ppu.oamData[i*4+2]
	var address uint16
	//if ppu.flagSpriteSize == 0 {
  if ppu.CTRL.SpriteSize == SPRITE_SIZE_8 {
		//if attributes&0x80 == 0x80 {
		//	row = 7 - row
		//}
		//table := ppu.flagSpriteTable
		//address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
    address = ppu.CTRL.SpritePatTableAddr + uint16(tile)*16 + uint16(row)
	//} else {
	//	//if attributes&0x80 == 0x80 {
	//	//	row = 15 - row
	//	//}
	//	//table := tile & 1
	//	tile &= 0xFE
	//	if row > 7 {
	//		tile++
	//		row -= 8
	//	}
	//	//address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
  //  address = ppu.CTRL.SpritePatTableAddr + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := ppu.Read(address)
	highTileByte := ppu.Read(address + 8)
  //log.Printf("Got tiles: %x, %x", ppu.Read(address), ppu.Read(address + 8))
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		//if attributes&0x40 == 0x40 {
		//	p1 = (lowTileByte & 1) << 0
		//	p2 = (highTileByte & 1) << 1
		//	lowTileByte >>= 1
		//	highTileByte >>= 1
		//} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		//}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (ppu *PPU) spritePixel2() (uint8, uint8) {
  if !ppu.MASK.showSprites {
    return 0, 0
  }

	for oamIdx := 0; oamIdx < ppu.spriteCount2; oamIdx++ {
    entry := ppu.oamEntries[oamIdx]


    // We know this sprite is visible in this scanline, but is it
    // visible in this column?
    if int(entry.x) <= ppu.Dot && int(entry.x) + 8 > ppu.Dot {
      fy := uint8((ppu.Dot) - int(entry.x))

      pixel := (
        (entry.attrs & 0x3) << 2                   |
        ((entry.tileHigh >> (7 - fy)) & 0x1) << 1  |
        ((entry.tileLow  >> (7 - fy)) & 0x1) << 0)

      if pixel % 4 == 0 {
        continue
      }

      index := uint8(oamIdx)
      return index, pixel
    }
  }

  return 0, 0
}

func (ppu *PPU) spritePixel3() (byte, byte) {
	if ppu.MASK.showSprites == false {
		return 0, 0
	}

	for i := 0; i < ppu.spriteCount2; i++ {
    entry := ppu.oamEntries[i]

		offset := (ppu.Dot - 1) - int(entry.x)
		if offset < 0 || offset > 7 {
			continue
		}
    //log.Printf("Willlll fy %x", offset)
		offset = 7 - offset
		//color := byte((ppu.spritePatterns[i] >> byte(offset*4)) & 0x0F)
    //  fy := uint8(offset)
    //  color = (
    //    (entry.attrs & 0x3) << 2                    |
    //    ((entry.tileHigh >> (15 - fy)) & 0x1) << 1 |
    //    ((entry.tileLow  >> (15 - fy)) & 0x1) << 0)
		//color := byte((entry.patt >> byte(offset*4)) & 0x0F)

    var color uint8
    fy := uint8(offset)

    color = ((entry.attrs & 0x3) << 2          |
        ((entry.tileHigh >> (fy)) & 0x1) << 1  |
        ((entry.tileLow  >> (fy)) & 0x1) << 0)

		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (ppu *PPU) spritePixel() (byte, byte) {
	if ppu.MASK.showSprites == false {
		return 0, 0
	}

	for i := 0; i < ppu.spriteCount; i++ {
		offset := (ppu.Dot - 1) - int(ppu.spritePositions[i])
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
