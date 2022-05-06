package collections

func ItemCollection(gameName, collectionName string) (result *CollectionInfo, e error) {
	// Check if collection and collection info exists
	e = FullCollectionCheck(gameName, collectionName)
	if e != nil {
		return
	}

	// Get info
	result, e = CollectionGetInfo(gameName, collectionName)
	if e != nil {
		return
	}
	return
}
