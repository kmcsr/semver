
package semver

import (
	"encoding/json"
	"errors"
	"strconv"
)

// A [vaild semver](https://semver.org/#backusnaur-form-grammar-for-valid-semver-versions) can be parsed to Version
// Version also can parse from some partially missing semver, to allow reuse when parsing version range
type Version struct {
	Major, Minor, Patch int
	Pre string
	Build string
}

var Nil = Version{
	Major: -1,
	Minor: -1,
	Patch: -1,
}

func parseRemain(s string)(v Version, remain string, err error){
	if s == "" {
		return Nil, "", nil
	}
	s, v.Build = split(s, '+')
	s, v.Pre = split(s, '-')

	var a string
	a, s = split(s, '.')
	if a == "" || a == "*" || a == "x" || a == "X" {
		v.Major = -1
	}else{
		if v.Major, err = strconv.Atoi(a); err != nil {
			return
		}
		a, s = split(s, '.')
		if a == "" || a == "*" || a == "x" || a == "X" {
			v.Minor = -1
		}else{
			if v.Minor, err = strconv.Atoi(a); err != nil {
				return
			}
			a, s = split(s, '.')
			if a == "" || a == "*" || a == "x" || a == "X" {
				v.Patch = -1
			}else{
				if v.Patch, err = strconv.Atoi(a); err != nil {
					return
				}
			}
		}
	}
	remain = s
	return
}

func Parse(s string)(v Version, err error){
	if v, s, err = parseRemain(s); err != nil {
		return
	}
	if len(s) > 0 {
		err = errors.New("Unexpected character " + strconv.QuoteRune((rune)(s[0])))
		return
	}
	return
}

// IsValid check if the version is a valid semver or not
// See: <https://semver.org/#backusnaur-form-grammar-for-valid-semver-versions>
func (v Version)IsValid()(bool){
	// TODO: is it necessary to check Pre and Build?
	return v.Major >= 0 && v.Minor >= 0 && v.Patch >= 0
}

func (v Version)String()(s string){
	if v.Major >= 0 {
		s = strconv.Itoa(v.Major)
		if v.Minor >= 0 {
			s += "." + strconv.Itoa(v.Minor)
			if v.Patch >= 0 {
				s += "." + strconv.Itoa(v.Patch)
			}
		}
	}
	if len(v.Pre) > 0 {
		s += "-" + v.Pre
	}
	if len(v.Build) > 0 {
		s += "+" + v.Build
	}
	return
}

func (v Version)Compare(o Version)(int){
	if v.Major < 0 || o.Major < 0 {
		return 0
	}
	if v.Major < o.Major {
		return -1
	}
	if v.Major > o.Major {
		return 1
	}
	if v.Minor < 0 || o.Minor < 0 {
		return 0
	}
	if v.Minor < o.Minor {
		return -1
	}
	if v.Minor > o.Minor {
		return 1
	}
	if v.Patch < 0 || o.Patch < 0 {
		return 0
	}
	if v.Patch < o.Patch {
		return -1
	}
	if v.Patch > o.Patch {
		return 1
	}
	if v.Pre < o.Pre {
		return -1
	}
	if v.Pre > o.Pre {
		return 1
	}
	if v.Build < o.Build {
		return -1
	}
	if v.Build > o.Build {
		return 1
	}
	return 0
}

var _ json.Unmarshaler = (*Version)(nil)
var _ json.Marshaler = (*Version)(nil)

func (v Version)MarshalJSON()([]byte, error){
	return json.Marshal(v.String())
}

func (v *Version)UnmarshalJSON(data []byte)(err error){
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	*v, err = Parse(s)
	return
}

var _ Comparator = (*Version)(nil)

func (v Version)Contains(o Version)(bool){
	return v.Compare(o) == 0
}
