package replace

import (
	"github.com/HardDie/DeckBuilder/internal/apperr"
	servicesReplace "github.com/HardDie/DeckBuilder/internal/services/replace"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type Replace struct {
	svc servicesReplace.Replace
}

func New(svc servicesReplace.Replace) *Replace {
	return &Replace{svc: svc}
}

type PrepareResult struct {
	Data []servicesReplace.Couple `json:"data"`
}

type Result struct {
	Data *tts_entity.RootObjects `json:"data"`
}

// Prepare lists unique FaceURL and BackURL keys from a Saved Object JSON file.
func (r *Replace) Prepare(file []byte) (*PrepareResult, error) {
	if len(file) == 0 {
		return nil, apperr.With(apperr.ErrBadRenderFile, "choose a render file (.json)")
	}
	resp, err := r.svc.Prepare(file)
	if err != nil {
		return nil, err
	}
	return &PrepareResult{Data: resp}, nil
}

// Replace rewrites FaceURL and BackURL using the mapping JSON the GUI used to post.
func (r *Replace) Replace(file, mapping []byte) (*Result, error) {
	if len(file) == 0 {
		return nil, apperr.With(apperr.ErrBadRenderFile, "choose a render file (.json)")
	}
	if len(mapping) == 0 {
		return nil, apperr.With(apperr.ErrBadMappingFile, "choose a mapping file (.json)")
	}
	resp, err := r.svc.Replace(file, mapping)
	if err != nil {
		return nil, err
	}
	return &Result{Data: resp}, nil
}
