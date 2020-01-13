// Tideland Go Text - Simple Markup Language
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package sml // import "tideland.dev/go/text/sml"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"tideland.dev/go/trace/failure"
)

//--------------------
// SML WRITER
//--------------------

// WriterProcessor a processor
type WriterProcessor interface {
	// WriterProcessor is a processor itself.
	Processor

	// SetContext sets the writer context.
	SetContext(ctx *WriterContext)
}

// WriterProcessors is a map of processors that can be plugged
// into the used SML writer. The key is the first part of each
// tag.
type WriterProcessors map[string]WriterProcessor

// WriterContext controls the parameters of a writing.
type WriterContext struct {
	plugins     WriterProcessors
	writer      io.Writer
	prettyPrint bool
	indentStr   string
}

// NewWriterContext creates a new writer context.
func NewWriterContext(wp WriterProcessor, w io.Writer, pretty bool, indentStr string) *WriterContext {
	ctx := &WriterContext{
		plugins:     WriterProcessors{"": wp},
		writer:      w,
		prettyPrint: pretty,
		indentStr:   indentStr,
	}
	wp.SetContext(ctx)
	return ctx
}

// Register adds a writer processor which will be responsible for
// processing if the tag is matching.
func (ctx *WriterContext) Register(tag string, wp WriterProcessor) error {
	if _, ok := ctx.plugins[tag]; ok {
		return failure.New("plugin '%s' already registered", tag)
	}
	wp.SetContext(ctx)
	ctx.plugins[tag] = wp
	return nil
}

// Writef writes a formatted string to the writer.
func (ctx *WriterContext) Writef(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(ctx.writer, format, args...)
	return err
}

// mlWriter writes it to an io.Writer.
type mlWriter struct {
	context *WriterContext
	stack   []WriterProcessor
	indent  int
}

// WriteSML uses one WriterProcessor and possible more as plugins to write the
// SML node to the passed writer. The prettyPrint flag controls if the writing
// is in a more beautiful formatted way.
func WriteSML(node Node, ctx *WriterContext) error {
	wp := ctx.plugins[""]
	if wp == nil {
		return failure.New("need root processor")
	}
	w := &mlWriter{
		context: ctx,
		stack:   []WriterProcessor{wp},
		indent:  0,
	}
	// Process the node with the new writer.
	return node.ProcessWith(w)
}

// OpenTag writes the opening of a tag.
func (w *mlWriter) OpenTag(tag []string) error {
	w.activatePlugin(tag[0])
	w.writeIndent(true)
	if err := w.activePlugin().OpenTag(tag); err != nil {
		return err
	}
	w.writeNewline()
	w.indent++
	return nil
}

// CloseTag writes the closing of a tag.
func (w *mlWriter) CloseTag(tag []string) error {
	w.indent--
	w.writeIndent(false)
	if err := w.activePlugin().CloseTag(tag); err != nil {
		return err
	}
	w.writeNewline()
	// Check if a plugin is to deactivate.
	if len(w.stack) > 1 {
		w.deactivatePlugin()
	}
	return nil
}

// Text writes a text with an encoding of special runes.
func (w *mlWriter) Text(text string) error {
	w.writeIndent(true)
	if err := w.activePlugin().Text(text); err != nil {
		return err
	}
	w.writeNewline()
	return nil
}

// Raw writes raw data without any encoding.
func (w *mlWriter) Raw(raw string) error {
	w.writeIndent(true)
	if err := w.activePlugin().Raw(raw); err != nil {
		return err
	}
	w.writeNewline()
	return nil
}

// Comment writes comment data without any encoding.
func (w *mlWriter) Comment(comment string) error {
	w.writeIndent(true)
	if err := w.activePlugin().Comment(comment); err != nil {
		return err
	}
	w.writeNewline()
	return nil
}

// activatePlugin activates a new one of the
// registered plugins.
func (w *mlWriter) activatePlugin(tag string) {
	if p := w.context.plugins[tag]; p != nil {
		w.stack = append(w.stack, p)
	}
}

// deactivatePlugin deactivates the top plugin.
func (w *mlWriter) deactivatePlugin() {
	w.stack = w.stack[:len(w.stack)-1]
}

// activePlugin returns the current active plugin.
func (w *mlWriter) activePlugin() WriterProcessor {
	return w.stack[len(w.stack)-1]
}

// writeIndent writes an indentation if wanted.
func (w *mlWriter) writeIndent(open bool) {
	if w.context.prettyPrint {
		for i := 0; i < w.indent; i++ {
			if err := w.context.Writef(w.context.indentStr); err != nil {
				panic(err)
			}
		}
	} else if open {
		if err := w.context.Writef(" "); err != nil {
			panic(err)
		}
	}
}

// writeNewline writes a newline if wanted.
func (w *mlWriter) writeNewline() {
	if w.context.prettyPrint {
		if err := w.context.Writef("\n"); err != nil {
			panic(err)
		}
	}
}

//--------------------
// STANDARD SML WRITER PROCESSOR
//--------------------

// standardSMLWriter writes a SML document in its standard
// notation to a writer.
type standardSMLWriter struct {
	context *WriterContext
}

// NewStandardSMLWriter creates a new writer for a ML
// document in standard notation.
func NewStandardSMLWriter() WriterProcessor {
	return &standardSMLWriter{}
}

// SetContext sets the writer context.
func (w *standardSMLWriter) SetContext(ctx *WriterContext) {
	w.context = ctx
}

// OpenTag writes the opening of a tag.
func (w *standardSMLWriter) OpenTag(tag []string) error {
	return w.context.Writef("{%s", strings.Join(tag, ":"))
}

// CloseTag writes the closing of a tag.
func (w *standardSMLWriter) CloseTag(tag []string) error {
	return w.context.Writef("}")
}

// Text writes a text with an encoding of special runes.
func (w *standardSMLWriter) Text(text string) error {
	var buf bytes.Buffer
	for _, r := range text {
		switch r {
		case '^':
			buf.WriteString("^^")
		case '{':
			buf.WriteString("^{")
		case '}':
			buf.WriteString("^}")
		default:
			buf.WriteRune(r)
		}
	}
	return w.context.Writef(buf.String())
}

// Raw writes raw data without any encoding.
func (w *standardSMLWriter) Raw(raw string) error {
	return w.context.Writef("{! %s !}", raw)
}

// Comment writes comment data without any encoding.
func (w *standardSMLWriter) Comment(comment string) error {
	return w.context.Writef("{# %s #}", comment)
}

//--------------------
// XML WRITER PROCESSOR
//--------------------

// xmlWriter writes a ML document in XML notation.
type xmlWriter struct {
	context *WriterContext
	rawTag  string
}

// NewXMLWriter creates a new writer for a ML
// document in XML notation.
func NewXMLWriter(rawTag string) WriterProcessor {
	return &xmlWriter{
		rawTag: rawTag,
	}
}

// SetContext sets the writer context.
func (w *xmlWriter) SetContext(ctx *WriterContext) {
	w.context = ctx
}

// OpenTag writes the opening of a tag.
func (w *xmlWriter) OpenTag(tag []string) error {
	if err := w.context.Writef("<%s", tag[0]); err != nil {
		return err
	}
	if len(tag) > 1 {
		if err := w.context.Writef(" id=%q", tag[1]); err != nil {
			return err
		}
	}
	if len(tag) > 2 {
		if err := w.context.Writef(" class=%q", tag[2]); err != nil {
			return err
		}
	}
	return w.context.Writef(">")
}

// CloseTag writes the closing of a tag.
func (w *xmlWriter) CloseTag(tag []string) error {
	return w.context.Writef("</%s>", tag[0])
}

// Text writes a text with an encoding of special runes.
func (w *xmlWriter) Text(text string) error {
	return xml.EscapeText(w.context.writer, []byte(text))
}

// Raw writes raw data without any encoding.
func (w *xmlWriter) Raw(raw string) error {
	return w.context.Writef("<%s>%s</%s>", w.rawTag, raw, w.rawTag)
}

// Comment writes comment data without any encoding.
func (w *xmlWriter) Comment(comment string) error {
	return w.context.Writef("<!-- %s -->", comment)
}

// EOF
