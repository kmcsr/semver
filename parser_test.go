
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
		{},
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

