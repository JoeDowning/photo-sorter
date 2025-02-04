package genutils

func PrefixZeros(reqLen int, s string) string {
	if len(s) >= reqLen {
		return s
	}
	return "0" + PrefixZeros(reqLen-1, s)
}
