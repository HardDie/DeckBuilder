// Measure pages: cell size, grid, and file paths. No pixels.
package layout

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"

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

// Pages splits each deck into sheets of at most 69 faces.
// The cell comes from the first face. Later pages keep that cell.
func Pages(dir string, decks []Deck, scale int, shadow bool) ([]Page, []File, error) {
	var pages []Page
	var backs []File
	var common int
	for _, deck := range decks {
		if len(deck.Faces) == 0 {
			continue
		}
		common++
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
			name := fmt.Sprintf("%d_%s_%d_%d_%dx%d.jpg", common, deck.ID, index, len(faces), cols, rows)
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
				common++
			}
			faces = append(faces, face)
		}
		if len(faces) > 0 {
			flush()
		}
	}
	return pages, backs, nil
}

func backName(dir, title string, raw []byte) string {
	sum := md5.Sum(raw)
	name := "backside_" + title + "_" + fmt.Sprintf("%x", sum[0:3]) + ".png"
	return fs.PathToAbsolutePath(filepath.Join(dir, name))
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
