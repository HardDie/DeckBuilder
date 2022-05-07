package collections

type CollectionInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func NewCollectionInfo(id, name, desc, image string) *CollectionInfo {
	return &CollectionInfo{
		Id:          id,
		Name:        name,
		Description: desc,
		Image:       image,
	}
}

func (i *CollectionInfo) Compare(val *CollectionInfo) bool {
	if i.Id != val.Id {
		return false
	}
	if i.Name != val.Name {
		return false
	}
	if i.Description != val.Description {
		return false
	}
	if i.Image != val.Image {
		return false
	}
	return true
}
