package sets

func DeduplicateIntSlice(items []int) []int {
	seen := make(map[int]struct{})
	result := make([]int, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
