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
	"errors"
	"strconv"
)

//--------------------
// CONSTANTS
//--------------------

// NodeType describe the JSON type of node, may it be the
// root or any lower one.
type NodeType int

const (
	NodeTypeNull NodeType = iota
	NodeTypeObject
	NodeTypeArray
	NodeTypeString
	NodeTypeNumber
	NodeTypeBool
)

//--------------------
// NODE HELPERS
//--------------------

// nodeType determines the JSON type of a node (root or value).
func nodeType(data interface{}) NodeType {
	if data == nil {
		return NodeTypeNull
	}
	switch data.(type) {
	case map[string]interface{}:
		return NodeTypeObject
	case []interface{}:
		return NodeTypeArray
	case string:
		return NodeTypeString
	case int, float64:
		return NodeTypeNumber
	case bool:
		return NodeTypeBool
	}
	panic("invalid node type")

}

// nodeLen returns the length of the passed data (which can be a single
// value too).
func nodeLen(data interface{}) int {
	if data == nil {
		return 0
	}
	switch d := data.(type) {
	case []interface{}:
		return len(d)
	case map[string]interface{}:
		return len(d)
	case string, int, float64, bool:
		return 1
	}
	return 0
}

// nodeDo performs a function on all elements of the passed node (which
// can be a single value too).
func nodeDo(path []string, data interface{}, f func(k string, v *Value) error) error {
	if data == nil {
		return nil
	}
	switch d := data.(type) {
	case []interface{}:
		for i, d := range d {
			k := "#" + strconv.Itoa(i)
			if err := f(k, newValue(path, d, nil)); err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		for k, d := range d {
			if err := f(k, newValue(path, d, nil)); err != nil {
				return err
			}
		}
		return nil
	case string, int, float64, bool:
		return f("", newValue(path, data, nil))
	}
	return nil
}

// nodeAt finds a node at a given path of keys. It works recursive and
// collects the already done path for a possible error.
func nodeAt(data interface{}, done, path []string) (interface{}, error) {
	if len(path) == 0 {
		return data, nil
	}
	switch d := data.(type) {
	case map[string]interface{}:
		value, ok := d[path[0]]
		if !ok {
			return nil, &PathError{
				Mode: "object",
				Path: append(done, path[0]),
				Err:  errors.New("path does not exist"),
			}
		}
		if len(path) > 1 {
			return nodeAt(value, append(done, path[0]), path[1:])
		}
		return value, nil
	case []interface{}:
		index, err := indexOf(path[0])
		if err != nil {
			return nil, &PathError{
				Mode: "array",
				Path: append(done, path[0]),
				Err:  err,
			}
		}
		if index < 0 || index > len(d)-1 {
			return nil, &PathError{
				Mode: "array",
				Path: append(done, path[0]),
				Err:  errors.New("invalid array index"),
			}
		}
		value := d[index]
		if len(path) > 1 {
			return nodeAt(value, append(done, path[0]), path[1:])
		}
		return value, nil
	default:
		return nil, &PathError{
			Mode: "value",
			Path: append(done, path...),
			Err:  errors.New("path too long"),
		}
	}
}

// indexOf tries to convert an index string like "#5" into an integer
// index like 5.
func indexOf(index string) (int, error) {
	if index[0] != '#' {
		return 0, errors.New("no index")
	}
	return strconv.Atoi(index[1:])
}

// EOF
