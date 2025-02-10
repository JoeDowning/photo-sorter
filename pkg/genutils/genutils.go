package genutils

import "strings"

func PrefixZeros(reqLen int, s string) string {
	if len(s) >= reqLen {
		return s
	}
	return "0" + PrefixZeros(reqLen-1, s)
}

func StringsContainInArray(arr []string, str string) bool {
	for _, a := range arr {
		if strings.Contains(str, a) {
			return true
		}
	}
	return false
}

func InArray(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func ContainsFromArray(arr []string, str string) bool {
	for _, a := range arr {
		if strings.Contains(str, a) {
			return true
		}
	}
	return false
}
