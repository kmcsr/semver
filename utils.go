
package semver

import (
	"strings"
)

func split(s string, b byte)(l, r string){
	i := strings.IndexByte(s, b)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i + 1:]
}
