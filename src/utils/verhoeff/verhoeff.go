package verhoeff

import "unicode"

var verhoeff_d = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
	{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
	{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
	{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
	{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
	{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
	{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
	{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
	{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
}

// The permutation table
var verhoeff_p = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
	{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
	{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
	{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
	{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
	{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
	{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
}

// The inverse table
var verhoeff_inv = []int{0, 4, 3, 2, 1, 5, 6, 7, 8, 9}

// GenerateVerhoeff will For a given number generates a Verhoeff digit
func GenerateVerhoeff(num string) int {
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = verhoeff_d[c][verhoeff_p[((i + 1) % 8)][num[ll-i-1]-'0']]
	}
	return verhoeff_inv[c]
}

// GenerateVerhoeffString adds a checksum as the last digit to a numeric string.
func GenerateVerhoeffString(s string) (newS string) {
	r := GenerateVerhoeff(s)
	newS = s + []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}[r%10]
	return
}

// ValidateVerhoeff returns true if the passed string 'num' is Verhoeff compliant.  The check digit must be the last one.
func ValidateVerhoeff(num string) bool {
	if !IsInt(num) {
		return false
	}
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = verhoeff_d[c][verhoeff_p[(i % 8)][num[ll-i-1]-'0']]
	}
	return (c == 0)
}

// ValidateAndStrip checks the string, if it is valid then the string w/o the checksum is returned.
func ValidateAndStrip(num string) (ok bool, s string) {
	if ValidateVerhoeff(num) {
		return true, num[0 : len(num)-1]
	}
	return false, ""
}

func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}