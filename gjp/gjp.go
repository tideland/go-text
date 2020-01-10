// Tideland Go Text - Generic JSON Processor
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp // import "tideland.dev/go/text/gjp"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"

	"tideland.dev/go/text/stringex"
	"tideland.dev/go/trace/failure"
)

//--------------------
// DOCUMENT
//--------------------

// PathValue is the combination of path and value.
type PathValue struct {
	Path  string
	Value *Value
}

// PathValues contains a number of path/value combinations.
type PathValues []PathValue

// ValueProcessor describes a function for the processing of
// values while iterating over a document.
type ValueProcessor func(path string, value *Value) error

// Document represents one JSON document.
type Document struct {
	separator string
	root      interface{}
}

// Parse reads a raw document and returns it as
// accessible document.
func Parse(data []byte, separator string) (*Document, error) {
	var root interface{}
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, failure.Annotate(err, "cannot unmarshal document")
	}
	return &Document{
		separator: separator,
		root:      root,
	}, nil
}

// NewDocument creates a new empty document.
func NewDocument(separator string) *Document {
	return &Document{
		separator: separator,
	}
}

// Length returns the number of elements for the given path.
func (d *Document) Length(path string) int {
	n, err := valueAt(d.root, splitPath(path, d.separator))
	if err != nil {
		return -1
	}
	// Check if object or array.
	o, ok := isObject(n)
	if ok {
		return len(o)
	}
	a, ok := isArray(n)
	if ok {
		return len(a)
	}
	return 1
}

// SetValueAt sets the value at the given path.
func (d *Document) SetValueAt(path string, value interface{}) error {
	parts := splitPath(path, d.separator)
	root, err := setValueAt(d.root, value, parts)
	if err != nil {
		return err
	}
	d.root = root
	return nil
}

// ValueAt returns the addressed value.
func (d *Document) ValueAt(path string) *Value {
	n, err := valueAt(d.root, splitPath(path, d.separator))
	return &Value{n, err}
}

// Clear removes the so far build document data.
func (d *Document) Clear() {
	d.root = nil
}

// Query allows to find pathes matching a given pattern.
func (d *Document) Query(pattern string) (PathValues, error) {
	pvs := PathValues{}
	err := d.Process(func(path string, value *Value) error {
		if stringex.Matches(pattern, path, false) {
			pvs = append(pvs, PathValue{
				Path:  path,
				Value: value,
			})
		}
		return nil
	})
	return pvs, err
}

// Process iterates over a document and processes its values.
// There's no order, so nesting into an embedded document or
// list may come earlier than higher level paths.
func (d *Document) Process(processor ValueProcessor) error {
	return process(d.root, []string{}, d.separator, processor)
}

// MarshalJSON implements json.Marshaler.
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.root)
}

// EOF
