// WebP encode and decode through libwebp (BSD), linked as a static archive.
// The Makefile passes the include folder and libwebp.a with libsharpyuv.a,
// so this file names no library: nothing can pick a shared libwebp. See ADR 027.
package webp

/*
#include <stdlib.h>
#include <webp/encode.h>
#include <webp/decode.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"unsafe"
)

// MaxSide is the largest width or height a WebP image can have.
const MaxSide = 16383

var (
	// ErrEncode is a failed encode; the cause follows it.
	ErrEncode = errors.New("webp encode")
	// ErrDecode is data that is not a readable WebP image.
	ErrDecode = errors.New("webp decode")
)

// IsWebP reports whether data starts with a WebP header ("RIFF", size, "WEBP").
func IsWebP(data []byte) bool {
	return len(data) >= 12 && bytes.Equal(data[0:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP"))
}

// Encode writes img as lossy WebP at quality 0–100. Alpha is kept;
// an image whose alpha is all 255 gets no alpha channel.
func Encode(img *image.NRGBA, quality float32) ([]byte, error) {
	return encode(img, func(pix *C.uint8_t, w, h, stride C.int, out **C.uint8_t) C.size_t {
		return C.WebPEncodeRGBA(pix, w, h, stride, C.float(quality), out)
	})
}

// EncodeLossless writes img as lossless WebP. Every visible pixel decodes exactly;
// the hidden color of a fully transparent pixel may change.
func EncodeLossless(img *image.NRGBA) ([]byte, error) {
	return encode(img, func(pix *C.uint8_t, w, h, stride C.int, out **C.uint8_t) C.size_t {
		return C.WebPEncodeLosslessRGBA(pix, w, h, stride, out)
	})
}

// encode checks the size, runs one libwebp simple-API encoder, and copies its output.
func encode(img *image.NRGBA, run func(pix *C.uint8_t, w, h, stride C.int, out **C.uint8_t) C.size_t) ([]byte, error) {
	w, h := img.Rect.Dx(), img.Rect.Dy()
	if w <= 0 || h <= 0 || w > MaxSide || h > MaxSide {
		return nil, fmt.Errorf("%w: size %dx%d, each side must be 1–%d", ErrEncode, w, h, MaxSide)
	}
	var out *C.uint8_t
	// The pixel buffer holds no Go pointers, so passing it to C is allowed.
	n := run((*C.uint8_t)(unsafe.Pointer(&img.Pix[img.PixOffset(img.Rect.Min.X, img.Rect.Min.Y)])),
		C.int(w), C.int(h), C.int(img.Stride), &out)
	if n == 0 {
		return nil, fmt.Errorf("%w: libwebp returned no data", ErrEncode)
	}
	defer C.WebPFree(unsafe.Pointer(out))
	return C.GoBytes(unsafe.Pointer(out), C.int(n)), nil
}

// Size reads width and height from the header without decoding.
func Size(data []byte) (int, int, error) {
	if len(data) == 0 {
		return 0, 0, ErrDecode
	}
	var w, h C.int
	if C.WebPGetInfo((*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &w, &h) == 0 {
		return 0, 0, ErrDecode
	}
	return int(w), int(h), nil
}

// Decode reads WebP into NRGBA (color not premultiplied by alpha).
func Decode(data []byte) (*image.NRGBA, error) {
	w, h, err := Size(data)
	if err != nil {
		return nil, err
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	ok := C.WebPDecodeRGBAInto(
		(*C.uint8_t)(unsafe.Pointer(&data[0])), C.size_t(len(data)),
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])), C.size_t(len(img.Pix)), C.int(img.Stride))
	if ok == nil {
		return nil, ErrDecode
	}
	return img, nil
}
