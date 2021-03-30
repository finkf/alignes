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
			if got, _ := alignLines(tc.gt, tc.ocr); !reflect.DeepEqual(got.tab, tc.want.tab) {
				t.Errorf("exepected %#v; got %#v", tc.want.tab, got.tab)
			}
		})
	}
}

func TestFoo(t *testing.T) {
	gt := []string{
		"R aris ſed pedibus placet altera, nam dolorem",
		"I ndocilis tolerare fugit: fugientem a vafra vo-",
		"luptas",
		"C aptat &amp; obleans ſummo de vertice montis",
		"V rget prcipitem Phlegethontis ad antra profunda.",
		"S ed qvibus e major Prudentia gnara viarum,",
		"A dſpiciunt Oculô vigilante Duceſ viaſ.",
		"M uſas proptereà udiô es probiore ſecutus",
		"M agne Vir ingeniô, oﬃciō vocante Deō:",
		"E ligat ut juum tua lea ſcientia caem,",
		"R eigio â qvali deſcendat origine: qvis t",
		"S ecurus vit cleia trames ad ara,",
		"B eatu ſacrrepetit facundia ſvad.",
		"A bt iniqva lues livoris turpe venenum,",
		"C onlio jui qvod tentat obee Laboris,",
		"H oc precor hoc voveo: mea vota ſecundet Jeſus.",
	}
	ocr := []string{
		"R ari ſed pedibus placet altera, namq. dolore",
		".e",
		"ndocilis tolerare fugit:",
		"ff va ⁊o-",
		"lupt a",
		"aptat ollectan ſummo vertice montis",
		"ee præcipit ont aantra profn.",
		"ed dena darr",
		"aiciunt OculvigDlanta ⸗.",
		"M uſas proiore ecutus",
		"M agne Vir neni oõ. docante Deã:",
		"ligat ut juſtum tua lecta ſcient callem,",
		"elligio a deſcenat originc: d ſit",
		"ccurus vitæ æle ect tc att",
		"e tuæ ſacrærepetit facundia ſvadæ.",
		"bſit inuiqva lues livor turpe dcm",
		"⸗",
	}
	want := "#d#d#############i"
	if _, got := alignLines(gt, ocr); got != want {
		t.Errorf("expected %s; got %s", want, got)
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		dir  string
		want bool
	}{
		{"testdata", true},
		{"nonexistent", false},
	}
	for _, tc := range tests {
		t.Run(tc.dir, func(t *testing.T) {
			if got := exists(tc.dir); got != tc.want {
				t.Errorf("expected %t; got %t", tc.want, got)
			}
		})
	}
}
