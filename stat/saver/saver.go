package saver

import (
	"fmt"

	"github.com/sergey-chebanov/fire/stat/record"
)

type Interface interface {
	Save([]*record.Record)
	Close()
}

type Constructor func(arguments string) (Interface, error)

var Constructors = make(map[string]Constructor)

func New(driver string, path string) (Interface, error) {
	if path == "" {
		return nil, nil
	}

	if constructor, ok := Constructors[driver]; ok {
		return constructor(path)
	}
	return nil, fmt.Errorf(`Can't find "%s"-type saver`, driver)
}
