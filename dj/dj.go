// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj // import "tideland.dev/go/text/dj"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"io"
)

//--------------------
// DOCUMENT
//--------------------

// Document represents one JSON document.
type Document struct {
	root interface{}
}

// New creates a new empty document.
func New() *Document {
	return &Document{}
}

// Parse reads a raw document from a reader and returns it as
// accessible document.
func Parse(r io.Reader) (*Document, error) {
	var bs []byte
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, &DocumentError{
			Action: "read document to parse",
			Err:    err,
		}
	}
	var root interface{}
	if err := json.Unmarshal(bs, &root); err != nil {
		return nil, &DocumentError{
			Action: "unmarshal document",
			Err:    err,
		}
	}
	return &Document{
		root: root,
	}, nil
}

// Len returns the number of elements on the root level of the document.
func (d *Document) Len() int {
	return nodeLen(d.root)
}

// At retrieves a value at a given path of keys.
func (d *Document) At(path ...string) *Value {
	data, err := nodeAt(d.root, []string{}, path)
	if err != nil {
		return newValue(path, nil, err)
	}
	return newValue(path, data, nil)
}

// EOF
