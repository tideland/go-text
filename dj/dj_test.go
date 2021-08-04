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

// TestParseDocument verifies the parsing of a document.
func TestParseDocument(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		name   string
		in     string
		length int
		err    string
	}{
		{
			"no document",
			``,
			0,
			"unexpected end of JSON input",
		}, {
			"empty document",
			`{}`,
			0,
			"",
		}, {
			"single string value",
			`"test"`,
			1,
			"",
		}, {
			"single integer value",
			`12345`,
			1,
			"",
		}, {
			"key/value document",
			`{"test": 12345}`,
			1,
			"",
		}, {
			"list document",
			`[1, 2, 3, 4, 5]`,
			5,
			"",
		}, {
			"nested document",
			`{"s": "string","l":[1,2,3],"r":[{"x":1,"y":2},{"x":2}]}`,
			3,
			"",
		}, {
			"invalid document (open list)",
			`{"s": [}`,
			0,
			"invalid character '}' looking for beginning of value",
		}, {
			"invalid document (open structure)",
			`{"s": {}`,
			0,
			"unexpected end of JSON input",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			defer assert.SetFailable(t)()
			b := bytes.NewBufferString(test.in)
			doc, err := dj.Parse(b)
			if test.err != "" {
				assert.ErrorContains(err, test.err)
			} else {
				assert.NoError(err)
				assert.NotNil(doc)
				assert.Length(doc.Root(), test.length)
			}
		})
	}
}

// TestDocumentAt verifies the navigation to a value of a document.
func TestDocumentAt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		name  string
		in    string
		path  []string
		value string
		err   string
	}{
		{
			"single string value",
			`"test"`,
			[]string{},
			"test",
			"",
		}, {
			"key/value document",
			`{"test": "12345"}`,
			[]string{"test"},
			"12345",
			"",
		}, {
			"list document",
			`["1", "2", "3", "4", "5"]`,
			[]string{"#2"},
			"3",
			"",
		}, {
			"nested document",
			`{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			[]string{"o", "a", "#3"},
			"4",
			"",
		}, {
			"not existing path",
			`{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			[]string{"o", "oops", "a"},
			"",
			"path does not exist",
		}, {
			"path too long",
			`{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			[]string{"o", "x", "oops"},
			"",
			"path too long",
		}, {
			"invalid index / no number",
			`{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			[]string{"o", "a", "oops"},
			"",
			"no index",
		}, {
			"invalid index / invalid number",
			`{"s": "string","o":{"x":"foo","a":["1","2","3","4","5"]}}`,
			[]string{"o", "a", "#999"},
			"",
			"invalid array index",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			defer assert.SetFailable(t)()
			b := bytes.NewBufferString(test.in)
			doc, err := dj.Parse(b)
			assert.NoError(err)
			value := doc.At(test.path...)
			if test.err != "" {
				assert.ErrorContains(value.Error(), test.err)
			} else {
				assert.Equal(value.String(), test.value)
			}
		})
	}
}

// EOF
