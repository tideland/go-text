// Tideland Go Text - Dynamic JSON - Testing
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj // import "tideland.dev/go/text/dj"

//--------------------
// TEST HELPER
//--------------------

// NewValue wraps the private value constructor for testing.
func NewValue(path []string, data interface{}, err error) *Value {
	return newValue(path, data, err)
}

// EOF
