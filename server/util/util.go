package util

func ArrayAdd(arr []string, newitem string) []string {
	c, _ := Contains(arr, newitem)
	if c {
		return arr

	}

	newarr := append(arr, newitem)
	return newarr
}

func ArrayRemove(arr []string, newitem string) []string {
	newarr := []string{}

	for _, item := range arr {
		if item == newitem {
			continue
		}

		newarr = append(newarr, item)
	}

	return newarr
}

func Contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, -1
}

func Handle(err error) {
	if err != nil {
		panic(err)
	}
}
