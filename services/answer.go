package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/innermond/printoo/models"
)

type Mistaker interface {
	HasErrors() bool
	GetErrors() Errors
	Errorless() interface{}
}

func ContextedMistaker(father context.Context, mk Mistaker) context.Context {
	ctx := context.TODO()
	if mk.HasErrors() {
		e := map[string]Errors{"errors": mk.GetErrors()}
		ctx = context.WithValue(father, ConvertJsonKey, e)
	} else {
		ctx = context.WithValue(father, ConvertJsonKey, mk.Errorless())
	}
	return ctx
}

type key int

const ConvertJsonKey = 0

type Errors map[string][]string

func (e Errors) ErrMessage(fld string, val interface{}) {
	tpl := "invalid value:"
	qv := "%q"
	if err, ok := val.(error); ok == true {
		val = err.Error()
		qv = "%v"
	}
	tpl += qv
	e[fld] = append(e[fld], fmt.Sprintf(tpl, val))
}

type Validation interface {
	Validate() bool
}

func CamelCase(k string) string {
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(k), "_", " ", -1)), " ", "", -1)
}

type storage struct{}

func (s storage) Init(dns string) error {
	return models.InitDatabase(dns)
}

func (s storage) End() {
	models.DB.Close()
}

var Storage storage
