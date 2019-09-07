package saver

import (
	"fmt"
	"strings"

	"github.com/sergey-chebanov/fire/stat/record"
)

type Interface interface {
	Save([]*record.Record)
	Close()
}

type Constructor func(arguments string) (Interface, error)

var Constructors = make(map[string]Constructor)

func New(path string) (Interface, error) {
	if path == "" {
		return nil, nil
	}

	pieces := strings.SplitN(path, ":", 2)
	if len(pieces) != 2 {
		return nil, fmt.Errorf("invalid path \"%s\" It should be in the form \"type:arguments\"", path)
	}
	if constructor, ok := Constructors[pieces[0]]; ok {
		return constructor(pieces[1])
	}
	return nil, fmt.Errorf(`Can't find "%s"-type saver`, pieces[0])
}
