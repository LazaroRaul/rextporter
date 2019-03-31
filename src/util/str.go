package util

// StrSliceContains return true if val found in items slice
func StrSliceContains(items []string, val string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}
	return false
}
