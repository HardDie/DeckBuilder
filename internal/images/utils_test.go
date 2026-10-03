package images

import (
	"bytes"
	"encoding/binary"
	stderrors "errors"
	"hash/crc32"
	"image"
	"image/png"
	"testing"

	"github.com/HardDie/DeckBuilder/internal/apperr"
)

// pngHeader returns a PNG that has only a signature and an IHDR chunk.
// DecodeConfig reads it; a full decode fails on the missing data.
func pngHeader(width, height uint32) []byte {
	var ihdr bytes.Buffer
	ihdr.WriteString("IHDR")
	_ = binary.Write(&ihdr, binary.BigEndian, width)
	_ = binary.Write(&ihdr, binary.BigEndian, height)
	ihdr.Write([]byte{8, 6, 0, 0, 0}) // 8-bit RGBA, no interlace

	var out bytes.Buffer
	out.WriteString("\x89PNG\r\n\x1a\n")
	_ = binary.Write(&out, binary.BigEndian, uint32(ihdr.Len()-4))
	out.Write(ihdr.Bytes())
	_ = binary.Write(&out, binary.BigEndian, crc32.ChecksumIEEE(ihdr.Bytes()))
	return out.Bytes()
}

func TestValidateImage(t *testing.T) {
	var small bytes.Buffer
	if err := png.Encode(&small, image.NewRGBA(image.Rect(0, 0, 100, 100))); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		input   []byte
		wantErr *apperr.Error
	}{
		{name: "small_png", input: small.Bytes()},
		{name: "not_an_image", input: []byte("text"), wantErr: apperr.ErrUnsupportedImage},
		// At the limit the header passes; decode then fails on the missing data.
		{name: "header_at_limit", input: pngHeader(16384, 8192), wantErr: apperr.ErrUnsupportedImage},
		{name: "header_over_limit", input: pngHeader(16384, 8193), wantErr: apperr.ErrImageTooLarge},
		{name: "bomb_header", input: pngHeader(50000, 50000), wantErr: apperr.ErrImageTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgType, err := ValidateImage(tt.input)
			if tt.wantErr == nil {
				if err != nil || imgType != "png" {
					t.Fatalf("got %q, %v", imgType, err)
				}
				return
			}
			if !stderrors.Is(err, tt.wantErr) {
				t.Fatalf("err %v, want %v", err, tt.wantErr)
			}
		})
	}
}
