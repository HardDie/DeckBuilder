package entity

import (
	"encoding/json"
	"time"

	"github.com/HardDie/DeckBuilder/internal/fsentry/internal/utils"
	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_error"
)

type Entry struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	CreatedAt *time.Time      `json:"createdAt"`
	UpdatedAt *time.Time      `json:"updatedAt"`
	Data      json.RawMessage `json:"data"`
}

func NewEntry(name string, data interface{}, isIndent bool) *Entry {
	var dataJson []byte
	if isIndent {
		dataJson, _ = json.MarshalIndent(data, "", "	")
	} else {
		dataJson, _ = json.Marshal(data)
	}
	return &Entry{
		Id:        utils.NameToID(name),
		Name:      name,
		CreatedAt: utils.Allocate(time.Now()),
		Data:      dataJson,
	}
}

func (i *Entry) SetName(name string) *Entry {
	i.Name = name
	i.Id = utils.NameToID(name)
	return i
}
func (i *Entry) UpdatedNow() *Entry {
	i.UpdatedAt = utils.Allocate(time.Now())
	return i
}
func (i *Entry) FlushTime() *Entry {
	i.CreatedAt = utils.Allocate(time.Now())
	i.UpdatedAt = nil
	return i
}
func (i *Entry) UpdateData(data interface{}, isIndent bool) error {
	var dataJson []byte
	var err error
	if isIndent {
		dataJson, err = json.MarshalIndent(data, "", "	")
	} else {
		dataJson, err = json.Marshal(data)
	}
	if err != nil {
		return fsentry_error.Wrap(err, fsentry_error.ErrorInternal)
	}
	i.Data = dataJson
	i.UpdatedAt = utils.Allocate(time.Now())
	return nil
}
