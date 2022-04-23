package collections

type CollectionInfoWithoutId struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type CollectionInfo struct {
	Id string `json:"id"`
	CollectionInfoWithoutId
}
