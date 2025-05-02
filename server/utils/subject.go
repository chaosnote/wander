package utils

import (
	"strings"
)

// params step
//
//	where<to/from>
//	action
//	who<opt>[uid|*|empty]
func Subject(params ...string) string {
	return strings.Join(params, ".")
}
