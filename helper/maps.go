package helper

import "sort"

// extracts string keys in ascending order
func Keys(m map[string]string) []string {
	result := make([]string, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}
	sort.Strings(result)
	return result
}
