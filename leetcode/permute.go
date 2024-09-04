package leetcode

func permute(str string) []string {
	var result []string
	backtrace([]byte(str), 0, &result)
	return result
}

func backtrace(str []byte, start int, result *[]string) {
	if start == len(str) {
		*result = append(*result, string(str))
		return
	}
	for i := start; i < len(str); i++ {
		swap(str, start, i)
		backtrace(str, start+1, result)
		swap(str, start, i)
	}
}

func swap(str []byte, i int, j int) {
	str[i], str[j] = str[j], str[i]
}
