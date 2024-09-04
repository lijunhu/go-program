package leetcode

func longestPalindrome(s string) string {

	length := len(s)
	if length <= 1 {
		return s
	}
	palindrome := ""

	for i := 0; i < length; i++ {
		odd := expendAroundCenter(s, i, i)
		if len(odd) > len(palindrome) {
			palindrome = odd
		}

		even := expendAroundCenter(s, i, i+1)
		if len(even) > len(palindrome) {
			palindrome = even
		}
	}
	return palindrome
}

func expendAroundCenter(s string, left, right int) string {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}
	return s[left+1 : right]
}
