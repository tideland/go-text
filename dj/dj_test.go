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

// TestNewDocument verifies the creation of an empty document.
func TestNewDocument(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	doc := dj.New()
	assert.NotNil(doc)
}

// TestParseDocument verifies the reading and parsing of an documents.
func TestParseDocument(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		description string
		in          string
		length      int
		err         string
	}{
		{
			description: "no document",
			in:          ``,
			err:         "unexpected end of JSON input",
		}, {
			description: "empty document",
			in:          `{}`,
			length:      0,
			err:         "",
		}, {
			description: "single string value",
			in:          `"test"`,
			length:      1,
			err:         "",
		}, {
			description: "single integer value",
			in:          `12345`,
			length:      1,
			err:         "",
		}, {
			description: "key/value document",
			in:          `{"test": 12345}`,
			length:      1,
			err:         "",
		}, {
			description: "list document",
			in:          `[1, 2, 3, 4, 5]`,
			length:      5,
			err:         "",
		}, {
			description: "nested document",
			in:          `{"s": "string","l":[1,2,3],"r":[{"x":1,"y":2},{"x":2}]}`,
			length:      3,
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
			assert.Length(doc, test.length)
		}
	}
}

// TestDocumentAt verifies the navigation to a value of a document.
func TestDocumentAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		description string
		in          string
		path        []string
		value       string
		err         string
	}{
		{
			description: "single string value",
			in:          `"test"`,
			path:        []string{},
			value:       "test",
		}, {
			description: "key/value document",
			in:          `{"test": "12345"}`,
			path:        []string{"test"},
			value:       "12345",
		}, {
			description: "list document",
			in:          `["1", "2", "3", "4", "5"]`,
			path:        []string{"#2"},
			value:       "3",
		}, {
			description: "nested document",
			in:          `{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			path:        []string{"o", "a", "#3"},
			value:       "4",
		}, {
			description: "not existing path",
			in:          `{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			path:        []string{"o", "oops", "a"},
			err:         "path does not exist",
		}, {
			description: "path too long",
			in:          `{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			path:        []string{"o", "oops", "a"},
			err:         "path too long",
		}, {
			description: "invalid index - no number",
			in:          `{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			path:        []string{"o", "a", "oops"},
			err:         "no index",
		}, {
			description: "invalid index - invalid number",
			in:          `{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			path:        []string{"o", "a", "#999"},
			err:         "invalid array index",
		},
	}
	for i, test := range tests {
		assert.Logf("running test #%d: %s", i, test.description)
		b := bytes.NewBufferString(test.in)
		doc, err := dj.Parse(b)
		assert.NoError(err)
		value := doc.At(test.path...)
		if test.err != "" {
			assert.ErrorContains(value.Error(), test.err)
		} else {
			assert.Equal(value.String(), test.value)
		}
	}
}

// EOF
