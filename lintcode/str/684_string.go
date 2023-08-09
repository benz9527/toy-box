package str

func MissingString(str1 string, str2 string) []string {
	// Write your code here
	spliteByBlank := func(str string) []string {
		oriList := make([]string, 0, 16)
		for {
			i, j := 0, 0
			for ; len(str) > j && str[j] != ' '; j++ {
			}
			if j >= len(str) {
				oriList = append(oriList, str)
				break
			}
			oriList = append(oriList, str[i:j])
			str = str[j+1:]
			if len(str) <= 0 {
				break
			}
		}
		return oriList
	}

	s1 := spliteByBlank(str1)
	s2 := spliteByBlank(str2)
	flags := make([]bool, len(s1))
	for _, s := range s2 {
		for i := 0; i < len(s1); i++ {
			if s1[i] == s {
				flags[i] = true
			}
		}
	}

	res := make([]string, 0, 8)
	for i, matched := range flags {
		if !matched {
			res = append(res, s1[i])
		}
	}
	return res
}
