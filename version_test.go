
package semver_test

import (
	"testing"
	. "github.com/kmcsr/semver"
)

type V = Version

func TestParseVersion(t *testing.T){
	type T struct {
		S string
		V V
		E bool
	}
	data := []T{
		// Valids
		{ "0.0.4", V{ 0, 0, 4, "", "" }, false },
		{ "1.2.3", V{ 1, 2, 3, "", "" }, false },
		{ "10.20.30", V{ 10, 20, 30, "", "" }, false },
		{ "1.1.2-prerelease+meta", V{ 1, 1, 2, "prerelease", "meta" }, false },
		{ "1.1.2+meta", V{ 1, 1, 2, "", "meta" }, false },
		{ "1.1.2+meta-valid", V{ 1, 1, 2, "", "meta-valid" }, false },
		{ "1.0.0-alpha", V{ 1, 0, 0, "alpha", "" }, false },
		{ "1.0.0-beta", V{ 1, 0, 0, "beta", "" }, false },
		{ "1.0.0-alpha.beta", V{ 1, 0, 0, "alpha.beta", "" }, false },
		{ "1.0.0-alpha.beta.1", V{ 1, 0, 0, "alpha.beta.1", "" }, false },
		{ "1.0.0-alpha.1", V{ 1, 0, 0, "alpha.1", "" }, false },
		{ "1.0.0-alpha0.valid", V{ 1, 0, 0, "alpha0.valid", "" }, false },
		{ "1.0.0-alpha.0valid", V{ 1, 0, 0, "alpha.0valid", "" }, false },
		{ "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", V{ 1, 0, 0, "alpha-a.b-c-somethinglong", "build.1-aef.1-its-okay" }, false },
		{ "1.0.0-rc.1+build.1", V{ 1, 0, 0, "rc.1", "build.1" }, false },
		{ "2.0.0-rc.1+build.123", V{ 2, 0, 0, "rc.1", "build.123" }, false },
		{ "1.2.3-beta", V{ 1, 2, 3, "beta", "" }, false },
		{ "10.2.3-DEV-SNAPSHOT", V{ 10, 2, 3, "DEV-SNAPSHOT", "" }, false },
		{ "1.2.3-SNAPSHOT-123", V{ 1, 2, 3, "SNAPSHOT-123", "" }, false },
		{ "1.0.0", V{ 1, 0, 0, "", "" }, false },
		{ "2.0.0", V{ 2, 0, 0, "", "" }, false },
		{ "1.1.7", V{ 1, 1, 7, "", "" }, false },
		{ "2.0.0+build.1848", V{ 2, 0, 0, "", "build.1848" }, false },
		{ "2.0.1-alpha.1227", V{ 2, 0, 1, "alpha.1227", "" }, false },
		{ "1.0.0-alpha+beta", V{ 1, 0, 0, "alpha", "beta" }, false },
		{ "1.2.3----RC-SNAPSHOT.12.9.1--.12+788", V{ 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", "788" }, false },
		{ "1.2.3----R-S.12.9.1--.12+meta", V{ 1, 2, 3, "---R-S.12.9.1--.12", "meta" }, false },
		{ "1.2.3----RC-SNAPSHOT.12.9.1--.12", V{ 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", "" }, false },
		{ "1.0.0+0.build.1-rc.10000aaa-kk-0.1", V{ 1, 0, 0, "", "0.build.1-rc.10000aaa-kk-0.1" }, false },
		{ "v1.2.3", V{ 1, 2, 3, "", "" }, false },
		{ "V1.2.3", V{ 1, 2, 3, "", "" }, false },
		// // Do not support big int yet
		// { "99999999999999999999999.999999999999999999.99999999999999999", V{ 99999999999999999999999, 999999999999999999, 99999999999999999, "", "" }, false },
		{ "1.0.0-0A.is.legal", V{ 1, 0, 0, "0A.is.legal", "" }, false },
		// Invalids
		{ "1", Nil, true },
		{ "1.2", Nil, true },
		// { "1.2.3-0123", Nil, true },
		// { "1.2.3-0123.0123", Nil, true },
		// { "1.1.2+.123", Nil, true },
		{ "+invalid", Nil, true },
		{ "-invalid", Nil, true },
		{ "-invalid+invalid", Nil, true },
		{ "-invalid.01", Nil, true },
		{ "alpha", Nil, true },
		{ "alpha.beta", Nil, true },
		{ "alpha.beta.1", Nil, true },
		{ "alpha.1", Nil, true },
		{ "alpha+beta", Nil, true },
		{ "alpha_beta", Nil, true },
		{ "alpha.", Nil, true },
		{ "alpha..", Nil, true },
		{ "beta", Nil, true },
		// { "1.0.0-alpha_beta", Nil, true },
		{ "-alpha.", Nil, true },
		// { "1.0.0-alpha..", Nil, true },
		// { "1.0.0-alpha..1", Nil, true },
		// { "1.0.0-alpha...1", Nil, true },
		// { "1.0.0-alpha....1", Nil, true },
		// { "1.0.0-alpha.....1", Nil, true },
		// { "1.0.0-alpha......1", Nil, true },
		// { "1.0.0-alpha.......1", Nil, true },
		// { "01.1.1", Nil, true },
		// { "1.01.1", Nil, true },
		// { "1.1.01", Nil, true },
		{ "1.2", Nil, true },
		{ "1.2.3.DEV", Nil, true },
		{ "1.2-SNAPSHOT", Nil, true },
		{ "1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788", Nil, true },
		{ "1.2-RC-SNAPSHOT", Nil, true },
		{ "-1.0.3-gamma+b7718", Nil, true },
		{ "+justmeta", Nil, true },
		// { "9.8.7+meta+meta", Nil, true },
		// { "9.8.7-whatever+meta+meta", Nil, true },
		{ "99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12", Nil, true },
	}
	for _, d := range data {
		v, e := Parse(d.S)
		if d.E {
			if e == nil && v.IsValid() {
				t.Errorf("Expect error when parsing version %q, but got %#v", d.S, v)
			}
		}else if e != nil || !v.IsValid() {
			t.Errorf("Unexpect error when parsing version %q: %v", d.S, e)
		}else if v.Compare(d.V) != 0 {
			t.Errorf("Unexpect version %#v when parsing version %q,\n  expect %#v", v, d.S, d.V)
		}
	}
}

