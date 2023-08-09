package str

const (
	startChar = '^'
	endChar   = '$'
	comma     = ','
)

func LongestPalindrome(s string) string {
	// Quick return
	if len(s) == 1 {
		return s
	}

	// Convert to extended string, even to odd
	extendedStr := string(startChar)
	for i, char := range s {
		extendedStr += string(char)
		if i != len(s) {
			extendedStr += string(comma)
		}
	}
	extendedStr += string(endChar)

	maxLong := ""
	for i, char := range extendedStr {
		if char == startChar {
			continue
		}
		if char == endChar {
			break
		}

		p0 := i
		lp, rp, radius := i-1, i+1, 0
		for {
			if lp < 0 || rp >= len(extendedStr) {
				break
			}
			if extendedStr[lp] == extendedStr[rp] && extendedStr[lp] == comma {
				lp -= 1
				rp += 1
				continue
			} else if extendedStr[lp] == extendedStr[rp] && isCommonChar(rune(extendedStr[lp])) {
				lp -= 1
				rp += 1
				radius += 1
			} else if extendedStr[lp] == startChar || extendedStr[rp] == endChar || extendedStr[lp] != extendedStr[rp] {
				break
			}
		}

		oriStr := ""
		if isCommonChar(char) {
			oriIndex := int((p0 - 1) / 2)
			oriStr = s[oriIndex-radius : oriIndex+radius+1]
		} else if char == comma {
			if radius <= 0 {
				continue
			}
			oriLeftMinIndex := int((p0 - (2*radius - 1) - 1) / 2)
			oriRightMaxIndex := int((p0 + (2*radius - 1) - 1) / 2)
			oriStr = s[oriLeftMinIndex : oriRightMaxIndex+1]
		}
		if len(oriStr) > len(maxLong) {
			maxLong = oriStr
		}
	}
	return maxLong
}

func isCommonChar(char rune) bool {
	return char != startChar && char != endChar && char != comma
}
