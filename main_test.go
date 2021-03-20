package main

import (
	"reflect"
	"testing"
)

func mkmat(r, c int, vals ...int) mat {
	m := newMat(r, c)
	if len(m.tab) != len(vals) {
		panic("bad values for matrix")
	}
	copy(m.tab, vals)
	return m
}

func mkstrs(strs ...string) []string {
	return strs
}

func TestAlignLines(t *testing.T) {
	for _, tc := range []struct {
		name    string
		gt, ocr []string
		want    mat
	}{
		{"test-1", mkstrs("testa", "testb"), mkstrs("testx", "testy"), mkmat(3, 3, 0, 1, 2, 1, 2, 3, 2, 3, 4)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := alignLines(tc.gt, tc.ocr); reflect.DeepEqual(got.tab, tc.want.tab) {
				t.Fatalf("exepected %#v; got %#v", tc.want.tab, got.tab)
			}
		})
	}
}
