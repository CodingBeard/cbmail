package cbmail

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Dependencies struct {
	Config Config
	Logger Logger
	ErrorHandler ErrorHandler
}

type Config interface {
	GetRequiredString(path string) (string, error)
}

type ErrorHandler interface {
	Error(e error)
	Recover()
}

type DefaultErrorHandler struct{}

func (d DefaultErrorHandler) Error(e error) {
	buf := make([]byte, 1000000)
	runtime.Stack(buf, false)
	buf = bytes.Trim(buf, "\x00")
	stack := string(buf)
	stackParts := strings.Split(stack, "\n")
	newStackParts := []string{stackParts[0]}
	newStackParts = append(newStackParts, stackParts[3:]...)
	stack = strings.Join(newStackParts, "\n")
	log.Println("ERROR", e.Error()+"\n"+stack)
}

func (d DefaultErrorHandler) Recover() {
	e := recover()

	if e != nil {
		err, ok := e.(error)

		if ok {
			d.Error(err)
		} else {
			d.Error(errors.New(fmt.Sprint(e)))
		}
	}
}

type Logger interface {
	InfoF(category string, message string, args ...interface{})
}

type DefaultLogger struct{}

func (d DefaultLogger) InfoF(category string, message string, args ...interface{}) {
	log.Println(category+":", fmt.Sprintf(message, args...))
}