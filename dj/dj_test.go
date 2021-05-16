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
		len         int
		err         string
	}{
		{
			description: "no document",
			in:          ``,
			err:         "unexpected end of JSON input",
		}, {
			description: "empty document",
			in:          `{}`,
			len:         0,
			err:         "",
		}, {
			description: "single string value",
			in:          `"test"`,
			len:         1,
			err:         "",
		}, {
			description: "single integer value",
			in:          `12345`,
			len:         1,
			err:         "",
		}, {
			description: "key/value document",
			in:          `{"test": 12345}`,
			len:         1,
			err:         "",
		}, {
			description: "list document",
			in:          `[1, 2, 3, 4, 5]`,
			len:         5,
			err:         "",
		}, {
			description: "nested document",
			in:          `{"s": "string","l":[1,2,3],"r":[{"x":1,"y":2},{"x":2}]}`,
			len:         3,
			err:         "",
		}, {
			description: "invalid document (open list)",
			in:          `{"s": [}`,
			err:         "invalid character '}' looking for beginning of value",
		}, {
			description: "invalid document (open structure)",
			in:          `{"s": {}`,
			err:         "unexpected end of JSON input",
		},
	}
	for i, test := range tests {
		assert.Logf("running test #%d: %s", i, test.description)
		b := bytes.NewBufferString(test.in)
		doc, err := dj.Parse(b)
		if test.err != "" {
			assert.ErrorContains(err, test.err)
		} else {
			assert.NoError(err)
			assert.NotNil(doc)
			assert.Length(doc, test.len)
		}
	}
}

// EOF
