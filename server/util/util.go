package util

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
