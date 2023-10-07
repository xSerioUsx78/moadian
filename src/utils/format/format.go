package format

import (
	"sort"
)

func MergeMaps(
	map1 map[string]any, 
	map2 map[string]string,
) map[string]any {
	mergedMap := make(map[string]any)
	for key, value := range map1 {
		mergedMap[key] = value
	}
	for key, value := range map2 {
		mergedMap[key] = value
	}
	return mergedMap
}

func GetSortedMapKeys(data map[string]any) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}