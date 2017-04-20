package chain

import (
	"net/http"
)

type Adapter func(http.Handler) http.Handler

func (c *Chain) Links(m ...Adapter) *Chain {
	c.Adapters = append(c.Adapters, m...)
	if c.Handler != nil {
		return c.Pull(c.Handler)
	}
	return c
}

func (c *Chain) Pull(h http.Handler) *Chain {
	m := c.Adapters
	x := len(m)
	for x > 0 {
		x--
		h = m[x](h)
	}
	c.Handler = h
	c.Adapters = nil
	return c
}

func New(h ...Adapter) *Chain {
	c := &Chain{}
	c.Adapters = nil
	return c.Links(h...)
}

type Chain struct {
	Adapters []Adapter
	Handler  http.Handler
}
