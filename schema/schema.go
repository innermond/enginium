package schema

import "github.com/gorilla/schema"

type Decoder interface {
	Decode(interface{}, map[string][]string) error
}

type decoder struct {
	// a decoder that can do real work
	real Decoder
}

func (dec *decoder) Decode(dst interface{}, src map[string][]string) error {
	return dec.real.Decode(dst, src)
}

func NewDecoder() *decoder {
	// the real decoder is from gorrila
	gd := schema.NewDecoder()
	return &decoder{gd}
}
