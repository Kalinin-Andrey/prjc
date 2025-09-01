package fasthttp_tools

import (
	"strconv"
	"strings"
)

func Uints2Str(uints *[]uint, sep *string) string {
	if uints == nil || len(*uints) == 0 {
		return ""
	}
	if sep == nil {
		s := ","
		sep = &s
	}
	strs := make([]string, 0, len(*uints))
	var i uint
	for _, i = range *uints {
		strs = append(strs, strconv.FormatUint(uint64(i), 10))
	}
	return strings.Join(strs, *sep)
}
