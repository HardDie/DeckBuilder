package page

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
)

type goldenDoc struct {
	BackFile string       `json:"backFile"`
	Pages    []goldenPage `json:"pages"`
}

type goldenPage struct {
	Index    int    `json:"index"`
	Faces    int    `json:"faces"`
	FaceFile string `json:"faceFile"`
	Columns  int    `json:"columns"`
	Rows     int    `json:"rows"`
}

func TestIntegrationGoldenBytes(t *testing.T) {
	inputs := readInputs(t)
	cases := []struct {
		name   string
		title  string
		scale  int
		common int
		shadow bool
		faces  [][]byte
	}{
		{name: "one", title: "one", scale: 1, common: 1, faces: [][]byte{inputs["face_a.png"]}},
		{name: "mix", title: "mix", scale: 1, common: 2, shadow: true, faces: [][]byte{inputs["face_a.png"], inputs["face_b.png"]}},
		{name: "half", title: "half", scale: 2, common: 4, faces: [][]byte{inputs["face_a.png"]}},
		{name: "wide", title: "wide", scale: 1, common: 7, faces: [][]byte{inputs["face_wide.png"]}},
		{name: "many", title: "many", scale: 1, common: 5, faces: manyFaces(inputs)},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			back := inputs["back.png"]
			if tt.name == "wide" {
				back = inputs["back_wide.png"]
			}
			got := drawDeck(t, tt.title, tt.scale, tt.common, tt.shadow, tt.faces, back)
			dir := filepath.Join("testdata", tt.name)
			compareDoc(t, dir, got.doc)
			for name, body := range got.files {
				want, err := os.ReadFile(filepath.Join(dir, name))
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(body, want) {
					t.Fatalf("%s: %d bytes, golden %d bytes", name, len(body), len(want))
				}
			}
		})
	}
}

type drawn struct {
	doc   goldenDoc
	files map[string][]byte
}

func drawDeck(t *testing.T, title string, scale, common int, shadow bool, faces [][]byte, back []byte) drawn {
	t.Helper()
	dir := t.TempDir()
	cfg := entitiesSettings.Default()
	cfg.EnableBackShadow = shadow
	cur := New(title, dir, scale, common, &cfg)
	var doc goldenDoc
	files := map[string][]byte{}

	for _, face := range faces {
		if cur.IsEmpty() {
			abs, err := cur.SetBacksideImageAndSave(back)
			if err != nil {
				t.Fatal(err)
			}
			doc.BackFile = filepath.Base(abs)
			body, err := os.ReadFile(abs)
			if err != nil {
				t.Fatal(err)
			}
			files[doc.BackFile] = body
		}
		if cur.IsFull() {
			doc.Pages = append(doc.Pages, savePage(t, cur, files))
			cur = (&Page{}).Inherit(cur)
		}
		if err := cur.AddImage(face); err != nil {
			t.Fatal(err)
		}
	}
	if !cur.IsEmpty() {
		doc.Pages = append(doc.Pages, savePage(t, cur, files))
	}
	return drawn{doc: doc, files: files}
}

func savePage(t *testing.T, cur *Page, files map[string][]byte) goldenPage {
	t.Helper()
	abs, cols, rows, err := cur.Save()
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(abs)
	if err != nil {
		t.Fatal(err)
	}
	name := filepath.Base(abs)
	files[name] = body
	return goldenPage{
		Index:    cur.GetIndex(),
		Faces:    cur.Size(),
		FaceFile: name,
		Columns:  cols,
		Rows:     rows,
	}
}

func compareDoc(t *testing.T, dir string, got goldenDoc) {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(got); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(dir, "sheet.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("sheet.json:\n%s\ngolden:\n%s", buf.Bytes(), want)
	}
}

func readInputs(t *testing.T) map[string][]byte {
	t.Helper()
	dir := filepath.Join("testdata", "input")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[entry.Name()] = body
	}
	return out
}

func manyFaces(inputs map[string][]byte) [][]byte {
	faces := make([][]byte, 0, 70)
	for i := 0; i < 69; i++ {
		faces = append(faces, inputs["face_a.png"])
	}
	faces = append(faces, inputs["face_b.png"])
	return faces
}
