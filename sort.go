package godiff

import "sort"

func sortPrimitive(value interface{}) interface{} {
	switch actual := value.(type) {
	case *[]int:
		clone := make([]int, len(*actual))
		copy(clone, *actual)
		sort.Ints(clone)
		return &clone
	case *[]string:
		clone := make([]string, len(*actual))
		copy(clone, *actual)
		sort.Strings(clone)
		return &clone
	case *[]float64:
		clone := make([]float64, len(*actual))
		copy(clone, *actual)
		sort.Float64s(clone)
		return &clone
	case *[]float32:
		clone := make([]float32, len(*actual))
		copy(clone, *actual)
		sort.SliceStable(clone, func(i, j int) bool {
			return clone[i] < clone[j]
		})
		return &clone
	}
	return value
}
