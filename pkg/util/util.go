package util

func TrimRepeatString(strs []string) []string {
	var newArr []string
	var strRecord = map[string]bool{}
	for _, v := range strs {
		if !strRecord[v] {
			newArr = append(newArr, v)
		}
		strRecord[v] = true
	}

	return newArr
}

func TrimRepeatInt(strs []int) []int {
	var newArr []int
	var intRecord = map[int]bool{}
	for _, v := range strs {
		if !intRecord[v] {
			newArr = append(newArr, v)
		}
		intRecord[v] = true
	}

	return newArr
}
