package services

type Hello interface {
	SayHello() string
}

type hello struct{}

func (h *hello) SayHello() string {
	return "Hello, world!"
}

func NewHello() Hello {
	return &hello{}
}
