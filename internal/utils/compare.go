package utils

func CompareMaps(one, two map[string]string) bool {
	// In case both maps are empty
	if len(one) == 0 && len(two) == 0 {
		return true
	}

	// If they have a different number of items
	if len(one) != len(two) {
		return false
	}

	// Compare each item
	for key, value := range one {
		cmp, ok := two[key]
		if !ok {
			// If a value with the same key does not exist
			return false
		}
		if value != cmp {
			// If the value is different
			return false
		}
	}

	return true
}
