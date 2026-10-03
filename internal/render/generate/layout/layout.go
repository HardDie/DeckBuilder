// Measure pages: cell size, grid, and file paths. No pixels.
package layout

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"

	"github.com/cespare/xxhash/v2"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/fs"
)

// Deck is one deck's raw image bytes, in card order.
type Deck struct {
	ID    string
	Back  []byte
	Faces [][]byte
}

// Page is one sheet the drawer can paint without measuring again.
type Page struct {
	DeckID    string
	Index     int
	Faces     [][]byte
	Back      []byte
	CellW     int
	CellH     int
	Cols      int
	Rows      int
	SheetPath string
	BackPath  string
	Shadow    bool
}

// File is a raw back PNG. The bytes are the catalog file, unchanged.
type File struct {
	Path string
	Body []byte
}

// sheetVersion is part of every page hash.
// Bump it when sheet drawing changes pixels (resize, paint, shadow, encoder, quality),
// so a page drawn by an older build is not reused.
const sheetVersion = 1

// Pages splits each deck into sheets of at most 69 faces.
// The cell comes from the first face. Later pages keep that cell.
// A sheet is named <deck>_<page>_<hash>.jpg. The hash covers everything drawn on it,
// so an unchanged page keeps its name across renders and can be reused.
func Pages(dir string, decks []Deck, scale int, shadow bool) ([]Page, []File, error) {
	var pages []Page
	var backs []File
	for _, deck := range decks {
		if len(deck.Faces) == 0 {
			continue
		}
		cellW, cellH, err := cell(deck.Faces[0], scale)
		if err != nil {
			return nil, nil, err
		}
		backPath := backName(dir, deck.ID, deck.Back)
		backs = append(backs, File{Path: backPath, Body: deck.Back})

		index := 1
		var faces [][]byte
		flush := func() {
			cols, rows := grid(len(faces) + 1)
			sum := pageHash(faces, deck.Back, cellW, cellH, cols, rows, shadow)
			name := fmt.Sprintf("%s_%d_%016x.jpg", deck.ID, index, sum)
			pages = append(pages, Page{
				DeckID:    deck.ID,
				Index:     index,
				Faces:     append([][]byte(nil), faces...),
				Back:      deck.Back,
				CellW:     cellW,
				CellH:     cellH,
				Cols:      cols,
				Rows:      rows,
				SheetPath: fs.PathToAbsolutePath(filepath.Join(dir, name)),
				BackPath:  backPath,
				Shadow:    shadow,
			})
		}
		for _, face := range deck.Faces {
			if len(faces) >= config.MaxCount {
				flush()
				faces = nil
				index++
			}
			faces = append(faces, face)
		}
		if len(faces) > 0 {
			flush()
		}
	}
	return pages, backs, nil
}

// backName is backside_<deck>_<hash>.png. The file is the original bytes.
func backName(dir, title string, raw []byte) string {
	name := fmt.Sprintf("backside_%s_%016x.png", title, xxhash.Sum64(raw))
	return fs.PathToAbsolutePath(filepath.Join(dir, name))
}

// pageHash covers every input that changes a sheet's pixels, in drawing order.
// Card names, descriptions, variables, and counts are not drawn, so they are not hashed.
// Every image is prefixed with its length, so two inputs cannot run together.
func pageHash(faces [][]byte, back []byte, cellW, cellH, cols, rows int, shadow bool) uint64 {
	h := xxhash.New()
	var buf [8]byte
	num := func(v int) {
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		_, _ = h.Write(buf[:])
	}
	blob := func(b []byte) {
		num(len(b))
		_, _ = h.Write(b)
	}
	shade := 0
	if shadow {
		shade = 1
	}
	for _, v := range []int{sheetVersion, cellW, cellH, cols, rows, shade, len(faces)} {
		num(v)
	}
	blob(back)
	for _, face := range faces {
		blob(face)
	}
	return h.Sum64()
}

type box struct {
	inner float64
	w     int
	h     int
	set   bool
}

// cell repeats the sheet fit rule on the face's header width and height.
func cell(raw []byte, scale int) (int, int, error) {
	w, h, err := headerSize(raw)
	if err != nil {
		return 0, 0, err
	}
	_, got := apply(w, h, scale, box{})
	return got.w, got.h, nil
}

func apply(cardW, cardH, scale int, cur box) (int, box) {
	inner := cur.inner
	if cardW*10 > 10_000 {
		inner = 10_000 / 10 / float64(cardW)
		if int(math.Trunc(float64(cardW)*inner)) > 10_000 {
			inner += 0.01
		}
	}
	if math.Trunc(float64(cardH)*inner)*7 > 10_000 {
		inner = 10_000 / 7 / float64(cardH)
		if int(math.Trunc(float64(cardH)*inner)) > 10_000 {
			inner += 0.01
		}
	}
	if scale == 0 {
		scale = 1
	}
	if inner == 0 {
		inner = 1
	}
	out := cur
	out.inner = inner
	if !cur.set {
		out.w = int(math.Trunc(float64(cardW)*inner)) / scale
		out.h = int(math.Trunc(float64(cardH)*inner)) / scale
		out.set = true
	}
	return scale, out
}

// grid is the smallest columns×rows that can hold n cells, from 2×2 to 10×7.
func grid(n int) (cols, rows int) {
	cols = 10
	rows = 7
	maxCards := cols * rows
	for r := config.MinHeight; r <= config.MaxHeight; r++ {
		for c := config.MinWidth; c <= config.MaxWidth; c++ {
			possible := c * r
			if possible < maxCards && possible >= n {
				maxCards = possible
				cols = c
				rows = r
			}
		}
	}
	return
}

func headerSize(raw []byte) (int, int, error) {
	if w, h, ok := pngSize(raw); ok {
		return w, h, nil
	}
	if w, h, ok := jpegSize(raw); ok {
		return w, h, nil
	}
	return 0, 0, fmt.Errorf("image size: not a png or jpeg")
}

func pngSize(raw []byte) (int, int, bool) {
	sig := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if len(raw) < 24 || string(raw[:8]) != string(sig) || string(raw[12:16]) != "IHDR" {
		return 0, 0, false
	}
	return int(binary.BigEndian.Uint32(raw[16:20])), int(binary.BigEndian.Uint32(raw[20:24])), true
}

func jpegSize(raw []byte) (int, int, bool) {
	if len(raw) < 4 || raw[0] != 0xff || raw[1] != 0xd8 {
		return 0, 0, false
	}
	i := 2
	for i+1 < len(raw) {
		if raw[i] != 0xff {
			return 0, 0, false
		}
		for i < len(raw) && raw[i] == 0xff {
			i++
		}
		if i >= len(raw) {
			return 0, 0, false
		}
		marker := raw[i]
		i++
		if marker == 0xd9 || marker == 0xda {
			return 0, 0, false
		}
		if marker >= 0xd0 && marker <= 0xd7 {
			continue
		}
		if i+1 >= len(raw) {
			return 0, 0, false
		}
		seg := int(binary.BigEndian.Uint16(raw[i : i+2]))
		if seg < 2 || i+seg > len(raw) {
			return 0, 0, false
		}
		if isSOF(marker) && seg >= 7 {
			h := int(binary.BigEndian.Uint16(raw[i+3 : i+5]))
			w := int(binary.BigEndian.Uint16(raw[i+5 : i+7]))
			return w, h, true
		}
		i += seg
	}
	return 0, 0, false
}

func isSOF(marker byte) bool {
	switch marker {
	case 0xc0, 0xc1, 0xc2, 0xc3, 0xc5, 0xc6, 0xc7, 0xc9, 0xca, 0xcb, 0xcd, 0xce, 0xcf:
		return true
	default:
		return false
	}
}
