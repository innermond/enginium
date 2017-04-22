package printoo

type AnswerOk struct {
	Code int
	Data interface{}
}

type AnswerBad struct {
	Code   int
	Errors error
}
