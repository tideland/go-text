// Tideland Go Text - Dynamic JSON - Testing
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj_test // import "tideland.dev/go/text/dj"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/text/dj"
)

//--------------------
// TESTS
//--------------------

// TestNew tests the creation of an empty document.
func TestNew(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	doc := dj.New()
	assert.NotNil(doc)
}

// TestParse tests the parsing of documents.
func TestParse(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		description string
		in          string
		expectDoc   bool
		err         string
	}{
		{
			description: "empty document",
			in: "",
			expectDoc: true,
			err: "",
		}
	}
	for i, test := range tests {
		assert.Logf("running test #%d: %s", i, test.description)
		br := bytes.NewBufferString(test.in)
		doc, err := dj.Parse(br)
		if test.expectDoc {
			assert.NotNil(doc)
		}
		if test.err != "" {
			assert.ErrorContains(err, test.err)
		}
	}
}

// EOF
