// Resize through stb_image_resize2 (v2.18, public domain / MIT, by Jeff Roberts
// and Jorge L Rodriguez). The header is compiled into the binary by cgo, so the
// app stays statically linked: there is no library to link against.
// SIMD is chosen by the compiler target: NEON on arm64, SSE2 on amd64.
// See ADR 022.
package stbresize

/*
#cgo CFLAGS: -O3
#define STB_IMAGE_RESIZE_IMPLEMENTATION
#define STB_IMAGE_RESIZE_STATIC
#include "stb_image_resize2.h"
*/
import "C"

import (
	"image"
	"image/draw"
	"unsafe"
)

// Resize scales img to width×height with the Catmull-Rom filter.
// Any decoded layout works (NRGBA, RGBA, YCbCr, Gray, Paletted, 16-bit);
// it is converted to 8-bit NRGBA first.
// An opaque image filters its 4 channels alike, which is about 2.4× faster.
// An image with transparency weights color by alpha,
// so hidden colors under transparent pixels do not bleed into edges.
func Resize(img image.Image, width, height int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	src := toNRGBA(img)
	if width <= 0 || height <= 0 || len(src.Pix) == 0 {
		return dst
	}

	layout := C.stbir_pixel_layout(C.STBIR_4CHANNEL)
	if !src.Opaque() {
		layout = C.STBIR_RGBA
	}
	// The buffers hold no Go pointers, so passing them to C is allowed.
	C.stbir_resize(
		unsafe.Pointer(&src.Pix[0]), C.int(src.Rect.Dx()), C.int(src.Rect.Dy()), C.int(src.Stride),
		unsafe.Pointer(&dst.Pix[0]), C.int(width), C.int(height), C.int(dst.Stride),
		layout, C.STBIR_TYPE_UINT8, C.STBIR_EDGE_CLAMP, C.STBIR_FILTER_CATMULLROM,
	)
	return dst
}

// toNRGBA returns img as an NRGBA whose Pix starts at its top-left pixel.
func toNRGBA(img image.Image) *image.NRGBA {
	b := img.Bounds()
	if src, ok := img.(*image.NRGBA); ok && b.Min == (image.Point{}) {
		return src
	}
	src := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(src, src.Bounds(), img, b.Min, draw.Src)
	return src
}
