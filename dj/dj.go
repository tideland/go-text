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

	"tideland.dev/go/trace/failure"
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
	if _, err := r.Read(bs); err != nil {
		return nil, failure.Annotate(err, "connot read document to parse")
	}
	var root interface{}
	if err := json.Unmarshal(bs, &root); err != nil {
		return nil, failure.Annotate(err, "cannot unmarshal document")
	}
	return &Document{
		root: root,
	}, nil
}

// EOF
