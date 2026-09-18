package entity

import (
	"encoding/json"
	"time"

	"github.com/HardDie/DeckBuilder/internal/fsentry/internal/utils"
	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_error"
	"github.com/HardDie/DeckBuilder/internal/fsentry/pkg/fsentry_types"
)

type FolderInfo struct {
	Id        string                     `json:"id"`
	Name      fsentry_types.QuotedString `json:"name"`
	CreatedAt *time.Time                 `json:"createdAt"`
	UpdatedAt *time.Time                 `json:"updatedAt"`
	Data      json.RawMessage            `json:"data"`
}

func NewFolderInfo(name string, data interface{}, isIndent bool) *FolderInfo {
	var dataJson []byte
	if isIndent {
		dataJson, _ = json.MarshalIndent(data, "", "	")
	} else {
		dataJson, _ = json.Marshal(data)
	}
	return &FolderInfo{
		Id:        utils.NameToID(name),
		Name:      fsentry_types.QuotedString(name),
		CreatedAt: utils.Allocate(time.Now()),
		Data:      dataJson,
	}
}

func (i *FolderInfo) SetName(name string) *FolderInfo {
	i.Name = fsentry_types.QuotedString(name)
	i.Id = utils.NameToID(name)
	return i
}
func (i *FolderInfo) UpdatedNow() *FolderInfo {
	i.UpdatedAt = utils.Allocate(time.Now())
	return i
}
func (i *FolderInfo) FlushTime() *FolderInfo {
	i.CreatedAt = utils.Allocate(time.Now())
	i.UpdatedAt = nil
	return i
}
func (i *FolderInfo) UpdateData(data interface{}, isIndent bool) error {
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
