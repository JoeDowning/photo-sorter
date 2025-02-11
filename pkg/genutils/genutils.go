package genutils

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

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

func OutputAllEntriesAtDir(logger *zap.Logger, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		logger.Fatal("failed to read directory", zap.Error(err))
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
