package document

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type Document interface {
	ValueMapper
}

type peekFunc func(r io.Reader) (Document, error)

type valueMap map[string]interface{}

type ValueMapper interface {
	ValueGetter
	ValueSetter
}

type ValueGetter interface {
	Get(path string) (interface{}, error)
}

type ValueSetter interface {
	Set(path string, v interface{}) error
}

func parse(r io.Reader, f ...peekFunc) (Document, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading document: %w", err)
	}

	for _, peek := range f {
		d, err := peek(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("error peeking document type: %w", err)
		}
		if d != nil {
			return d, nil
		}
	}

	return nil, fmt.Errorf("no valid document type")
}

func Parse(r io.Reader) (Document, error) {
	return parse(r, NewArgoCD)
}
