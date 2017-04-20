package domino

import "net/http"

type piece func(http.Handler) http.Handler

type domino struct {
	pieces []piece
	h      http.Handler
}

func Pieces(p ...piece) *domino {
	d := &domino{}
	d.pieces = p
	return d
}

func (d *domino) Roll(h http.Handler) http.Handler {
	x := len(d.pieces)
	for x > 0 {
		x--
		h = d.pieces[x](h)
	}
	return h
}
