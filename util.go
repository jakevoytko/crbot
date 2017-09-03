package main

// CollapseWhitespace returns an array where splitContent entries are removed,
// starting at the given index. Will return the input if no collapsing is
// necessary, otherwise will return a new slice with the given indices cut out.
func CollapseWhitespace(splitContent []string, startIndex int) []string {
	if startIndex >= len(splitContent) || len(splitContent[startIndex]) > 0 {
		return splitContent
	}

	endIndex := startIndex + 1

	for endIndex < len(splitContent) && len(splitContent[endIndex]) == 0 {
		endIndex++
	}

	return append(splitContent[:startIndex], splitContent[endIndex:]...)
}
