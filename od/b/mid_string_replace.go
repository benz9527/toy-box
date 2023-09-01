package b

import (
	isort "sort"
	"strings"
)

func GetReplaceString(str string) string {
	containers := [][]string{}
	inBrackets := false
	src := ""
	for i := 0; i < len(str); i++ {
		if str[i:i+1] == "(" {
			inBrackets = true
			containers = append(containers, []string{})
			continue
		} else if str[i:i+1] == ")" {
			inBrackets = false
			continue
		}

		if inBrackets {
			containers[len(containers)-1] = append(containers[len(containers)-1], str[i:i+1])
		} else {
			src += str[i : i+1]
		}
	}

	// Merge
	sliceToMap := func(s []string) map[string]struct{} {
		m := map[string]struct{}{}
		for _, item := range s {
			m[item] = struct{}{}
		}
		return m
	}
	in := func(cur map[string]struct{}, next []string) bool {
		for _, item := range next {
			if _, ok := cur[strings.ToLower(item)]; ok {
				return ok
			}
			if _, ok := cur[strings.ToUpper(item)]; ok {
				return ok
			}
		}
		return false
	}
	for i := 0; i < len(containers); i++ {
		curC := containers[i]
		for j := i + 1; j < len(containers); j++ {
			nextC := containers[j]
			if in(sliceToMap(curC), nextC) {
				containers[i] = append(curC, containers[i+1]...)
				containers[j] = []string{}
			}
		}
	}

	compares := []map[string]struct{}{}
	minChars := []string{}
	for i := 0; i < len(containers); i++ {
		if len(containers[i]) > 1 {
			compares = append(compares, sliceToMap(containers[i]))
			isort.Strings(containers[i])
			minChars = append(minChars, containers[i][0])
		}
	}

	ans := ""
	for j := 0; j < len(src); j++ {
		s := src[j : j+1]
		for i := 0; i < len(compares); i++ {
			if _, ok := compares[i][strings.ToLower(s)]; ok {
				s = minChars[i]
				continue
			}
			if _, ok := compares[i][strings.ToUpper(s)]; ok {
				s = minChars[i]
				continue
			}
		}
		ans += s
	}
	return ans
}
