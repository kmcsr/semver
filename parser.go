
package semver

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var _ json.Marshaler = (*ComparatorSet)(nil)
var _ json.Unmarshaler = (*ComparatorSet)(nil)

func ParseComparatorSet(s string)(rs ComparatorSet, err error){
	l := make(ComparatorList, 0, 3)
	for {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			if len(l) != 0 {
				rs = append(rs, l)
			}else if len(rs) != 0 {
				err = io.EOF
			}
			return
		}
		s0 := s[0]
		if len(s) >= 2 && s0 == '|' && s[1] == '|' {
			s = s[2:]
			rs = append(rs, l)
			l = make(ComparatorList, 0, 3)
			continue
		}else if '0' <= s0 && s0 <= '9' || s0 == '*' || s0 == 'x' || s0 == 'X' {
			var v Version
			var s2 string
			s2, s = split(s, ' ')
			if v, err = Parse(s2); err != nil {
				return
			}
			s = strings.TrimSpace(s)
			if len(s) > 0 && s[0] == '-' {
				s = strings.TrimSpace(s[1:])
				if len(s) == 0 {
					err = io.EOF
					return
				}
				r := &Range{
					Min: v,
				}
				var s2 string
				s2, s = split(s, ' ')
				if r.Max, err = Parse(s2); err != nil {
					return
				}
				l = append(l, r)
			}else{
				l = append(l, v)
			}
			continue
		}
		if len(s) == 1 {
			err = io.EOF
			return
		}
		var r Requirement
		switch s0 {
		case '=':
			if s[1] == '=' {
				s = s[2:]
			}else{
				s = s[1:]
			}
			r.Cond = EQ
		case '!':
			if s[1] == '=' {
				s = s[2:]
			}else{
				s = s[1:]
			}
			r.Cond = NE
		case '^':
			s = s[1:]
			r.Cond = EX
		case '~':
			s = s[1:]
			r.Cond = TD
		case '<':
			if s[1] == '=' {
				s = s[2:]
				r.Cond = LE
			}else{
				s = s[1:]
				r.Cond = LT
			}
		case '>':
			if s[1] == '=' {
				s = s[2:]
				r.Cond = GE
			}else{
				s = s[1:]
				r.Cond = GT
			}
		default:
			err = fmt.Errorf("Unexpected character %q", s[1])
			return
		}
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			err = io.EOF
			return
		}
		var s2 string
		s2, s = split(s, ' ')
		if r.Version, err = Parse(s2); err != nil {
			return
		}
		l = append(l, &r)
	}
}

func (rs *ComparatorSet)MarshalJSON()([]byte, error){
	return json.Marshal(rs.String())
}

func (rs *ComparatorSet)UnmarshalJSON(data []byte)(err error){
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	*rs, err = ParseComparatorSet(s)
	return
}
