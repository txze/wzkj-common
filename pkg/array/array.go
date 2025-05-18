package array

func IsInArrayString(strs []string, val string) bool {
	for _, v := range strs {
		if v == val {
			return true
		}
	}
	return false
}

func IsInArrayInt(values []int, val int) bool {
	for _, v := range values {
		if v == val {
			return true
		}
	}
	return false
}

func TrimRepeatString(values []string) []string {
	var newArr []string
	var record = map[string]bool{}
	for _, v := range values {
		if _, ok := record[v]; ok {
			continue
		}

		record[v] = true
		newArr = append(newArr, v)
	}

	return newArr
}
