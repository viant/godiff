package godiff

// Comparator an interface for comparison customization
type Comparator interface {
	Matches(from, to interface{}, tag *Tag) (bool, error)
}

func matches(from, to interface{}) bool {
	if from == nil {
		if to == nil {
			return true
		}
		return false
	}
	if to == nil {
		return false
	}
	return from == to
}
