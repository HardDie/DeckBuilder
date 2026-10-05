package images

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"testing"
)

// app1 is a JPEG APP1 segment whose EXIF IFD0 has one Orientation entry.
func app1(order binary.ByteOrder, orientation uint16) []byte {
	var tiff bytes.Buffer
	if order == binary.LittleEndian {
		tiff.WriteString("II")
	} else {
		tiff.WriteString("MM")
	}
	_ = binary.Write(&tiff, order, uint16(42))
	_ = binary.Write(&tiff, order, uint32(8)) // IFD0 right after the header
	_ = binary.Write(&tiff, order, uint16(1)) // one entry
	_ = binary.Write(&tiff, order, uint16(0x0112))
	_ = binary.Write(&tiff, order, uint16(3)) // SHORT
	_ = binary.Write(&tiff, order, uint32(1))
	_ = binary.Write(&tiff, order, orientation)
	_ = binary.Write(&tiff, order, uint16(0)) // pad the 4-byte value field
	_ = binary.Write(&tiff, order, uint32(0)) // no next IFD
	payload := append([]byte("Exif\x00\x00"), tiff.Bytes()...)
	seg := []byte{0xFF, 0xE1}
	seg = binary.BigEndian.AppendUint16(seg, uint16(len(payload)+2))
	return append(seg, payload...)
}

// withAPP1 puts an APP1 segment right after the JPEG start marker.
func withAPP1(jpg, seg []byte) []byte {
	out := append([]byte{}, jpg[:2]...)
	out = append(out, seg...)
	return append(out, jpg[2:]...)
}

func TestJPEGOrientation(t *testing.T) {
	soiEOI := []byte{0xFF, 0xD8, 0xFF, 0xD9}
	tests := []struct {
		name string
		data []byte
		want int
	}{
		{"little endian 6", withAPP1(soiEOI, app1(binary.LittleEndian, 6)), 6},
		{"big endian 3", withAPP1(soiEOI, app1(binary.BigEndian, 3)), 3},
		{"big endian 8", withAPP1(soiEOI, app1(binary.BigEndian, 8)), 8},
		{"no exif", soiEOI, 1},
		{"value out of range", withAPP1(soiEOI, app1(binary.LittleEndian, 9)), 1},
		{"not a jpeg", []byte("\x89PNG\r\n\x1a\n"), 1},
		{"truncated segment", []byte{0xFF, 0xD8, 0xFF, 0xE1, 0x10, 0x00, 'E'}, 1},
		{"other APP1 (XMP)", withAPP1(soiEOI, append([]byte{0xFF, 0xE1, 0x00, 0x08}, "http"...)), 1},
		{"empty", nil, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jpegOrientation(tt.data); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// TestOrient follows the source's top-left pixel (red) through each orientation.
func TestOrient(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 3, 2))
	src.SetNRGBA(0, 0, color.NRGBA{R: 255, A: 255})
	tests := []struct {
		orientation int
		size        image.Point
		red         image.Point
	}{
		{1, image.Pt(3, 2), image.Pt(0, 0)},
		{2, image.Pt(3, 2), image.Pt(2, 0)},
		{3, image.Pt(3, 2), image.Pt(2, 1)},
		{4, image.Pt(3, 2), image.Pt(0, 1)},
		{5, image.Pt(2, 3), image.Pt(0, 0)},
		{6, image.Pt(2, 3), image.Pt(1, 0)},
		{7, image.Pt(2, 3), image.Pt(1, 2)},
		{8, image.Pt(2, 3), image.Pt(0, 2)},
		{9, image.Pt(3, 2), image.Pt(0, 0)},
	}
	for _, tt := range tests {
		got := orient(src, tt.orientation)
		if got.Rect.Size() != tt.size {
			t.Errorf("orientation %d: size %v, want %v", tt.orientation, got.Rect.Size(), tt.size)
			continue
		}
		for y := 0; y < tt.size.Y; y++ {
			for x := 0; x < tt.size.X; x++ {
				isRed := got.NRGBAAt(x, y).R == 255
				if isRed != (image.Pt(x, y) == tt.red) {
					t.Errorf("orientation %d: pixel (%d,%d) red=%v, want red only at %v", tt.orientation, x, y, isRed, tt.red)
				}
			}
		}
	}
}
