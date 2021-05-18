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
	"fmt"
)

//--------------------
// ERRORS
//--------------------

// DocumentError records an error on higher document level.
type DocumentError struct {
	Action string
	Err    error
}

// Error represents the error as string.
func (de *DocumentError) Error() string {
	return fmt.Sprintf("%s: %v", de.Action, de.Err)
}

// Unwrap returns the internal error.
func (de *DocumentError) Unwrap() error {
	return de.Err
}

// PathError records an error when navigating inside a document.
type PathError struct {
	Mode string
	Path []string
	Err  error
}

// Error represents the error as string.
func (pe *PathError) Error() string {
	return fmt.Sprintf("%s at %v: %v", pe.Mode, pe.Path, pe.Err)
}

// Unwrap returns the internal error.
func (pe *PathError) Unwrap() error {
	return pe.Err
}

// EOF
