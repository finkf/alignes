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

func TestAlignLines(t *testing.T) {
	for _, tc := range []struct {
		name    string
		gt, ocr []string
		want    mat
	}{
		{"test-1", []string{"testa", "testb"}, []string{"testx", "testy"}, mkmat(3, 3, 0, 5, 10, 5, 1, 6, 10, 6, 2)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := alignLines(tc.gt, tc.ocr); !reflect.DeepEqual(got.tab, tc.want.tab) {
				t.Errorf("exepected %#v; got %#v", tc.want.tab, got.tab)
			}
		})
	}
}
