package parser

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}
