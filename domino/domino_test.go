package domino

import (
	"fmt"
	"net/http"
	"testing"
)

func TestInitialisation(t *testing.T) {
	d := new(domino)
	if d == nil {
		t.Error("domino init: ", d)
	}
}

func TestAddPieces(t *testing.T) {

	h1 := func(h http.Handler) http.Handler { return h }
	h2 := func(h http.Handler) http.Handler { return h }
	h3 := func(h http.Handler) http.Handler { return h }

	tests := []struct {
		hd []piece
		l  int
	}{
		{[]piece{h1}, 1},
		{[]piece{h1, h2}, 2},
		{[]piece{h1, h2, h3}, 3},
	}

	for _, test := range tests {
		d := Pieces(test.hd...)
		if len(d.pieces) != test.l {
			t.Errorf("expected %d got %d", test.l, len(d.pieces))
		}

	}

}

type h struct{}

func (hh *h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func TestRoll(t *testing.T) {

	d := new(domino)
	// seed handler
	h0 := new(h)
	// we got an improved handler
	h8 := d.Roll(h0)

	expect := "*domino.h"
	got := fmt.Sprintf("%T", h8)
	if got != expect {
		t.Errorf("expected %s got %s", got, expect)
	}
}
