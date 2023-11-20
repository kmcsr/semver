
package semver

import (
	"fmt"
	"strings"
)

type Comparator interface {
	Contains(Version)(bool)
	fmt.Stringer
}

// ComparatorList will pass when all of the ranges passed
type ComparatorList []Comparator

var _ Comparator = (*ComparatorList)(nil)

func (rs ComparatorList)Contains(v Version)(bool){
	for _, r := range rs {
		if !r.Contains(v) {
			return false
		}
	}
	return true
}

func (rs ComparatorList)String()(string){
	var b strings.Builder
	for i, r := range rs {
		_, isset := r.(*ComparatorSet)
		s := r.String()
		if i > 0 {
			b.WriteByte(' ')
		}
		if isset {
			b.WriteByte('(')
		}
		b.WriteString(s)
		if isset {
			b.WriteByte(')')
		}
	}
	return b.String()
}

// ComparatorSet will pass if any of the ranges passed
type ComparatorSet []Comparator

var _ Comparator = (*ComparatorSet)(nil)

func (rs ComparatorSet)Contains(v Version)(bool){
	if len(rs) == 0 {
		return true
	}
	for _, r := range rs {
		if r.Contains(v) {
			return true
		}
	}
	return false
}

func (rs ComparatorSet)String()(string){
	var b strings.Builder
	for i, r := range rs {
		s := r.String()
		if i > 0 {
			b.WriteString(" || ")
		}
		b.WriteString(s)
	}
	return b.String()
}

type Cond int
const (
	NoCond Cond = iota
	EQ // ==
	NE // !=
	LT // <
	GT // >
	LE // <=
	GE // >=
	EX // ^ ; same as 
	TD // ~
)

func CondFromString(s string)(c Cond){
	switch s {
	case "=", "==":
		return EQ
	case "!", "!=":
		return NE
	case "<":
		return LT
	case ">":
		return GT
	case "<=":
		return LE
	case ">=":
		return GE
	case "^":
		return EX
	case "~":
		return TD
	default:
		return NoCond
	}
}

func (c Cond)String()(string){
	switch c {
	case EQ: return "=="
	case NE: return "!="
	case LT: return "<"
	case GT: return ">"
	case LE: return "<="
	case GE: return ">="
	case EX: return "^"
	case TD: return "~"
	default: panic("Unknown cond")
	}
}

type Requirement struct {
	Cond Cond
	Version Version
}

var _ Comparator = (*Requirement)(nil)

func (r *Requirement)Contains(v Version)(bool){
	n := v.Compare(r.Version)
	switch r.Cond {
	case EQ:
		return n == 0
	case NE:
		return n != 0
	case LT:
		return n < 0
	case LE:
		return n <= 0
	case GT:
		return n > 0
	case GE:
		return n >= 0
	case EX:
		if r.Version.Major < 0 {
			return true
		}
		if r.Version.Major == 0 {
			if v.Major != 0 {
				return false
			}
			return r.Version.Minor < 0 || r.Version.Minor == v.Minor
		}
		return r.Version.Major == v.Major
	case TD:
		return r.Version.Major < 0 || r.Version.Major == v.Major && (
			r.Version.Minor < 0 || r.Version.Minor == v.Minor)
	default:
		panic("Unknown cond")
	}
}

func (r *Requirement)String()(string){
	return r.Cond.String() + r.Version.String()
}

type Range struct {
	Min, Max Version
}

var _ Comparator = (*Range)(nil)

func (r *Range)Contains(v Version)(bool){
	return r.Min.Compare(v) <= 0 && r.Max.Compare(v) >= 0
}

func (r *Range)String()(string){
	return r.Min.String() + " - " + r.Max.String()
}
