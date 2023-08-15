package ratelimtercore

// Filter function iterates over the slice and applies the filterFunc to each item
func filter(slice []string, filterFunc func(string) bool) []string {
	filtered := []string{}
	for _, item := range slice {
		if filterFunc(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
