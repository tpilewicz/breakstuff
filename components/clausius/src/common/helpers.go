package common

func SplitSSlice(slice []string, chunkSize int) [][]string {
	result := [][]string{}
	nbChunks := len(slice) / chunkSize
	if len(slice)%chunkSize != 0 {
		nbChunks += 1
	}
	for i := 0; i < nbChunks; i++ {
		result = append(result, []string{})
	}
	for i, elt := range slice {
		result[i/chunkSize] = append(result[i/chunkSize], elt)
	}
	return result
}

func SplitISlice(slice []int, chunkSize int) [][]int {
	result := [][]int{}
	nbChunks := len(slice) / chunkSize
	if len(slice)%chunkSize != 0 {
		nbChunks += 1
	}
	for i := 0; i < nbChunks; i++ {
		result = append(result, []int{})
	}
	for i, elt := range slice {
		result[i/chunkSize] = append(result[i/chunkSize], elt)
	}
	return result
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func KeysOfMap(m map[string]int) []string {
	result := []string{}
	for k := range m {
		result = append(result, k)
	}
	return result
}

func Except(s1 []string, s2 []string) []string {
	result := []string{}
	for _, k := range s1 {
		if !stringInSlice(k, s2) {
			result = append(result, k)
		}
	}
	return result
}
