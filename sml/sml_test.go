// Tideland Go Text - Simple Markup Language - Unit Tests
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package sml_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/text/sml"
)

//--------------------
// TESTS
//--------------------

// TestTagValidation checks if only correct tags are accepted.
func TestTagValidation(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	tests := []struct {
		in  string
		out []string
		ok  bool
	}{
		{"-abc", nil, false},
		{"-", nil, false},
		{"abc-", nil, false},
		{"ab-c", []string{"ab-c"}, true},
		{"abc", []string{"abc"}, true},
		{"ab:cd", []string{"ab", "cd"}, true},
		{"1a", nil, false},
		{"a1", []string{"a1"}, true},
		{"a:1", []string{"a", "1"}, true},
		{"a-b:c-d", []string{"a-b", "c-d"}, true},
		{"a-:c-d", nil, false},
		{"-a:c-d", nil, false},
		{"ab:-c", nil, false},
		{"ab:c-", nil, false},
		{"a-b-1", []string{"a-b-1"}, true},
		{"a-b-1:c-d-2:e-f-3", []string{"a-b-1", "c-d-2", "e-f-3"}, true},
	}
	for i, test := range tests {
		msg := fmt.Sprintf("%q (test %d) ", test.in, i)
		tag, err := sml.ValidateTag(test.in)
		if err == nil {
			assert.Equal(tag, test.out, msg)
			assert.True(test.ok, msg)
		} else {
			assert.ErrorMatch(err, fmt.Sprintf("invalid tag: %q", test.in), msg)
			assert.False(test.ok, msg)
		}
	}
}

// TestCreating checks the manual tree creation.
func TestCreating(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	root := createNodeStructure(assert)
	assert.Equal(root.Tag(), []string{"root"}, "Root tag has to be 'root'.")
	assert.NotEmpty(root, "Root tag is not empty.")
}

// TestWriterProcessing checks the writing of SML.
func TestWriterProcessing(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	root := createNodeStructure(assert)
	bufA := bytes.NewBufferString("")
	bufB := bytes.NewBufferString("")
	ctxA := sml.NewWriterContext(sml.NewStandardSMLWriter(), bufA, true, "    ")
	ctxB := sml.NewWriterContext(sml.NewStandardSMLWriter(), bufB, false, "")

	assert.NoError(sml.WriteSML(root, ctxA))
	assert.NoError(sml.WriteSML(root, ctxB))

	assert.Logf("===== WITH INDENT =====")
	assert.Logf(bufA.String())
	assert.Logf("===== WITHOUT INDENT =====")
	assert.Logf(bufB.String())
	assert.Logf("===== DONE =====")

	assert.NotEmpty(bufA, "Buffer A must not be empty.")
	assert.NotEmpty(bufB, "Buffer B must not be empty.")
}

// TestPositiveNodeReading checks the successful reading of nodes.
func TestPositiveNodeReading(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	text := "Before!   {foo:main {bar:1:first Yadda ^{Test^} 1} {! Raw: }} { ! ^^^ !}  {between}  {bar:2:last Yadda {Test ^^} 2}}   After!"
	builder := sml.NewNodeBuilder()
	err := sml.ReadSML(strings.NewReader(text), builder)
	assert.Nil(err)
	root, err := builder.Root()
	assert.Nil(err)
	assert.Equal(root.Tag(), []string{"foo", "main"})
	assert.NotEmpty(root)

	buf := bytes.NewBufferString("")
	ctx := sml.NewWriterContext(sml.NewStandardSMLWriter(), buf, true, "    ")
	assert.NoError(sml.WriteSML(root, ctx))

	assert.Logf("===== PARSED SML =====")
	assert.Logf(buf.String())
	assert.Logf("===== DONE =====")
}

// TestNegativeNodeReading checks the failing reading of nodes.
func TestNegativeNodeReading(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	text := "{Foo {bar:1 Yadda {test} {} 1} {bar:2 Yadda 2}}"
	builder := sml.NewNodeBuilder()
	err := sml.ReadSML(strings.NewReader(text), builder)
	assert.ErrorMatch(err, `.* invalid character after opening at index .*`)
}

// TestPositiveTreeReading checks the successful reading of trees.
func TestPositiveTreeReading(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	text := "{config {foo 1}{bar 2}{yadda {up down}{down up}}}"
	builder := sml.NewKeyStringValueTreeBuilder()
	err := sml.ReadSML(strings.NewReader(text), builder)
	assert.Nil(err)
	tree, err := builder.Tree()
	assert.Nil(err)
	assert.Logf("%v", tree)
}

// TestNegativeTreeReading checks the failing reading of trees.
func TestNegativeTreeReading(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	text := "{foo {bar 1}{bar 2}}"
	builder := sml.NewKeyStringValueTreeBuilder()
	err := sml.ReadSML(strings.NewReader(text), builder)
	assert.ErrorMatch(err, `.* node has multiple values`)
}

// TestSML2XML checks the conversion from SML to XML.
func TestSML2XML(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	in := `{html
{head {title A test document}}
{body
  {h1:title A test document}
  {p:intro:preface The is a simple sentence with an {em emphasized}
  and a {strong strong} text. We'll see how it renders.}
  {ul
    {li:1 It should be nice.}
    {li:2 It should be error free.}
    {li:3 It should be fast.}
  }
  {!
for foo := 0; foo < 42; foo++ {
	println(foo)
}
  !}
}}`
	builder := sml.NewNodeBuilder()
	err := sml.ReadSML(strings.NewReader(in), builder)
	assert.Nil(err)
	root, err := builder.Root()
	assert.Nil(err)

	buf := bytes.NewBufferString("")
	ctx := sml.NewWriterContext(sml.NewXMLWriter("pre"), buf, true, "    ")
	assert.NoError(ctx.Register("li", newLIWriter()))
	assert.NoError(sml.WriteSML(root, ctx))

	assert.Logf("===== XML =====")
	assert.Logf(buf.String())
	assert.Logf("===== DONE =====")
}

//--------------------
// HELPERS
//--------------------

// Create a node structure.
func createNodeStructure(assert *asserts.Asserts) sml.Node {
	check := assert.NoError
	builder := sml.NewNodeBuilder()

	check(builder.BeginTagNode("root"))

	check(builder.TextNode("Text A"))
	check(builder.TextNode("Text B"))
	check(builder.CommentNode("A first comment."))

	check(builder.BeginTagNode("sub-a:1st:important"))
	check(builder.TextNode("Text A.A"))
	check(builder.CommentNode("A second comment."))
	check(builder.EndTagNode())

	check(builder.BeginTagNode("sub-b:2nd"))
	check(builder.TextNode("Text B.A"))
	check(builder.BeginTagNode("text"))
	check(builder.TextNode("Any text with the special characters {, }, and ^."))
	check(builder.EndTagNode())
	check(builder.EndTagNode())

	check(builder.BeginTagNode("sub-c"))
	check(builder.TextNode("Before raw."))
	check(builder.RawNode("func Test(i int) { println(i) }"))
	check(builder.TextNode("After raw."))
	check(builder.EndTagNode())

	check(builder.EndTagNode())

	root, err := builder.Root()
	assert.Nil(err)

	return root
}

// liWriter handles the li-tag of the document.
type liWriter struct {
	context *sml.WriterContext
}

// newLIWriter creates a new writer for the li-tag.
func newLIWriter() sml.WriterProcessor {
	return &liWriter{}
}

// SetContext sets the writer context.
func (w *liWriter) SetContext(ctx *sml.WriterContext) {
	w.context = ctx
}

// OpenTag writes the opening of a tag.
func (w *liWriter) OpenTag(tag []string) error {
	return w.context.Writef("<li>")
}

// CloseTag writes the closing of a tag.
func (w *liWriter) CloseTag(tag []string) error {
	return w.context.Writef("</li>")
}

// Text writes a text with an encoding of special runes.
func (w *liWriter) Text(text string) error {
	return w.context.Writef("<em> %s </em>", text)
}

// Raw writes raw data without any encoding.
func (w *liWriter) Raw(raw string) error {
	return w.context.Writef("\n%s\n", raw)
}

// Comment writes comment data without any encoding.
func (w *liWriter) Comment(comment string) error {
	return w.context.Writef("\n%s\n", comment)
}

// EOF
