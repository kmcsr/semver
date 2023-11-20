
package semver_test

import (
	"testing"
	. "github.com/kmcsr/semver"
)


func TestParseComparator(t *testing.T){
	type T struct {
		S string
		E bool
	}
	data := []T{
		{"1.2.3", false},
		{"=1.2.3", false},
		{"==1.2.3", false},
		{"!1.2.3", false},
		{"!=1.2.3", false},
		{"<1.2.3", false},
		{">1.2.3", false},
		{"<=1.2.3", false},
		{">=1.2.3", false},
		{"~1.2.3", false},
		{"^1.2.3", false},
		{"1.2.3 ||", true},
		{"1.2.3 || ", true},
		{"1.2.3 || 2", false},
	}
	for _, d := range data {
		v, e := ParseComparatorSet(d.S)
		if d.E {
			if e == nil {
				t.Errorf("Expect error when parsing ComparatorSet %q, but got %#v", d.S, v)
			}
		}else if e != nil {
			t.Errorf("Unexpect error when parsing ComparatorSet %q: %v", d.S, e)
		}
	}
}

