package images

import (
	"bytes"
	"encoding/binary"
	"image"
)

// jpegOrientation reads the EXIF Orientation tag (1–8) of a JPEG.
// Anything missing or malformed is 1: the image is stored upright.
// Only the markers before the image data are read.
func jpegOrientation(data []byte) int {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return 1
	}
	for i := 2; i+4 <= len(data); {
		if data[i] != 0xFF {
			return 1
		}
		marker := data[i+1]
		if marker == 0xD8 || (marker >= 0xD0 && marker <= 0xD7) || marker == 0x01 || marker == 0xFF {
			i += 2
			continue
		}
		if marker == 0xDA || marker == 0xD9 { // start of scan or end: no EXIF before it
			return 1
		}
		size := int(binary.BigEndian.Uint16(data[i+2:]))
		end := i + 2 + size
		if size < 2 || end > len(data) {
			return 1
		}
		if marker == 0xE1 {
			if o := exifOrientation(data[i+4 : end]); o != 0 {
				return o
			}
		}
		i = end
	}
	return 1
}

// exifOrientation reads tag 0x0112 from IFD0 of an APP1 "Exif" payload. 0 means not found.
func exifOrientation(app1 []byte) int {
	if !bytes.HasPrefix(app1, []byte("Exif\x00\x00")) {
		return 0
	}
	tiff := app1[6:]
	if len(tiff) < 8 {
		return 0
	}
	var order binary.ByteOrder
	switch string(tiff[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 0
	}
	if order.Uint16(tiff[2:]) != 42 {
		return 0
	}
	ifd := int(order.Uint32(tiff[4:]))
	if ifd < 8 || ifd+2 > len(tiff) {
		return 0
	}
	count := int(order.Uint16(tiff[ifd:]))
	for n := 0; n < count; n++ {
		e := ifd + 2 + n*12
		if e+12 > len(tiff) {
			return 0
		}
		if order.Uint16(tiff[e:]) != 0x0112 {
			continue
		}
		// Type SHORT (3), count 1: the value sits in the first two bytes of the value field.
		if order.Uint16(tiff[e+2:]) != 3 {
			return 0
		}
		v := int(order.Uint16(tiff[e+8:]))
		if v < 1 || v > 8 {
			return 0
		}
		return v
	}
	return 0
}

// orient turns img so it stands upright, by EXIF orientation 1–8.
// 1 (or anything unknown) returns img unchanged. 5–8 swap width and height.
func orient(img *image.NRGBA, orientation int) *image.NRGBA {
	if orientation < 2 || orientation > 8 {
		return img
	}
	b := img.Rect
	w, h := b.Dx(), b.Dy()
	dw, dh := w, h
	if orientation >= 5 {
		dw, dh = h, w
	}
	dst := image.NewNRGBA(image.Rect(0, 0, dw, dh))
	for y := 0; y < dh; y++ {
		for x := 0; x < dw; x++ {
			var sx, sy int
			switch orientation {
			case 2: // mirrored left to right
				sx, sy = w-1-x, y
			case 3: // turned 180°
				sx, sy = w-1-x, h-1-y
			case 4: // mirrored top to bottom
				sx, sy = x, h-1-y
			case 5: // mirrored along the main diagonal
				sx, sy = y, x
			case 6: // needs a 90° turn clockwise
				sx, sy = y, h-1-x
			case 7: // mirrored along the other diagonal
				sx, sy = w-1-y, h-1-x
			case 8: // needs a 90° turn counterclockwise
				sx, sy = w-1-y, x
			}
			s := img.PixOffset(b.Min.X+sx, b.Min.Y+sy)
			d := dst.PixOffset(x, y)
			copy(dst.Pix[d:d+4], img.Pix[s:s+4])
		}
	}
	return dst
}
