package utils

func Find[T comparable](arr *[]T, el *T) int {
	for index, a := range *arr {
		if a == *el {
			return index
		}
	}
	return -1
}
