package validation

import (
	"regexp"
	"strings"
)

func IsPhone(p string, l int) bool {
	re := regexp.MustCompile("[0-9]+")
	nb := re.FindAllString(p, -1)
	s := []string{}
	for _, n := range nb {
		s = append(s, n)
	}
	phone := strings.Join(s, "")
	out := len(phone) >= l
	return out
}

func MinMax(min, max, val int, eq ...bool) bool {
	is := val < min || val > max
	if len(eq) > 0 && eq[0] {
		is = val <= min || val >= max

	}
	return is
}
