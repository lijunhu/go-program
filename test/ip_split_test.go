package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestSplitIp(t *testing.T) {

	var ipArray []string
	ipStr := "17.50.80.0-17.63.255.255"
	ipStrArray := strings.Split(ipStr, "-")
	a := strings.Split(ipStrArray[0], ".")
	b := strings.Split(ipStrArray[1], ".")
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			continue
		}


		tmp_a := PadLeft(a[i], 3-len(a[i]), "0")
		tmp_b := PadLeft(b[i], 3-len(b[i]), "0")

		int_a, _ := strconv.Atoi(tmp_a)
		int_b, _ := strconv.Atoi(tmp_b)

		for j := int_a; j <= int_b; j++ {
			switch i {
			case 0:
				ipArray = append(ipArray, fmt.Sprintf("%s.%s.%s.%s", string(j), a[i+1], a[i+2], a[i+3]))
			case 1:
				ipArray = append(ipArray, fmt.Sprintf("%s.%s.%s.%s", a[i-1], string(j), a[i+1], a[i+2]))
			case 2:
				ipArray = append(ipArray, fmt.Sprintf("%s.%s.%s.%s", a[i-2], a[i-1], string(j), a[i+1]))
			case 3:
				ipArray = append(ipArray, fmt.Sprintf("%s.%s.%s.%s", a[i-3], a[i-2], a[i-1], string(j)))
			}

		}

	}
	fmt.Print(len(ipArray))
}

func PadLeft(str string, width int, pad string) string {
	return Tile(pad, width-len(str)) + str
}

func PadRight(str string, width int, pad string) string {
	return str + Tile(pad, width-len(str))
}

func Pad(str string, width int, leftPad string, rightPad string) string {
	switch {
	case len(leftPad) == 0:
		return PadRight(str, width, rightPad)
	case len(rightPad) == 0:
		return PadLeft(str, width, leftPad)
	}
	padLen := (width - len(str)) / 2
	return Tile(leftPad, padLen) + str + Tile(rightPad, width-len(str)-padLen)
}

func Tile(pattern string, length int) string {
	patLen := len(pattern)
	if len(pattern) == 0 || length <= 0 {
		return ""
	}

	var buff strings.Builder
	for i := 0; i < length; i += patLen {
		buff.WriteString(pattern)
	}
	return MustSubstring(buff.String(), 0, length)
}

func MustSubstring(str string, start int, end int) string {
	res, err := Substring(str, start, end)
	if err != nil {
		panic(err)
	}
	return res
}

func Substring(str string, start int, end int) (string, error) {
	if start < 0 || start >= len(str) {
		return "", fmt.Errorf("start (%d) is out of range", start)
	}
	if end != 0 && end <= start {
		return "", fmt.Errorf("end (%d) cannot be equal to or smaller than start (%d)", end, start)
	}
	if end > len(str) {
		return "", fmt.Errorf("end (%d) is out of range", end)
	}

	var startByte = -1
	var runeIndex int
	for i := range str {
		if runeIndex == start {
			startByte = i
			if end == 0 {
				return str[startByte:], nil
			}
		}
		if end != 0 && runeIndex == end {
			return str[startByte:i], nil
		}
		runeIndex++
	}

	if startByte < 0 {
		return "", fmt.Errorf("start (%d) is out of range (%d)", start, runeIndex)
	}

	if end == runeIndex {
		return str[startByte:], nil
	}

	return "", fmt.Errorf("end (%d) is out of range (%d)", end, runeIndex)
}
