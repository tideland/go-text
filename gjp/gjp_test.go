// Tideland Go Text - Generic JSON Processor - Unit Tests
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp_test

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/text/gjp"
)

//--------------------
// TESTS
//--------------------

// TestParseError tests the returned error in case of
// an invalid document.
func TestParseError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs := []byte(`abc{def`)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(doc)
	assert.ErrorMatch(err, `.*cannot unmarshal document.*`)
}

// TestClear tests to clear a document.
func TestClear(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	doc.Clear()
	err = doc.SetValueAt("/", "foo")
	assert.NoError(err)
	foo := doc.ValueAt("/").AsString("<undefined>")
	assert.Equal(foo, "foo")
}

// TestLength tests retrieving values as strings.
func TestLength(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	l := doc.Length("X")
	assert.Equal(l, -1)
	l = doc.Length("")
	assert.Equal(l, 4)
	l = doc.Length("B")
	assert.Equal(l, 3)
	l = doc.Length("B/2")
	assert.Equal(l, 5)
	l = doc.Length("/B/2/D")
	assert.Equal(l, 2)
	l = doc.Length("/B/1/S")
	assert.Equal(l, 3)
	l = doc.Length("/B/1/S/0")
	assert.Equal(l, 1)
}

// TestProcessing tests the processing of adocument.
func TestProcessing(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)
	count := 0
	processor := func(path string, value *gjp.Value) error {
		count++
		assert.Logf("path %02d  =>  %-10q = %q", count, path, value.AsString("<undefined>"))
		return nil
	}

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	err = doc.Process(processor)
	assert.Nil(err)
	assert.Equal(count, 27)

	processor = func(path string, value *gjp.Value) error {
		return errors.New("ouch")
	}
	err = doc.Process(processor)
	assert.ErrorMatch(err, `.*ouch.*`)
}

// TestSeparator tests using different separators.
func TestSeparator(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, lo := createDocument(assert)

	// Slash as separator, once even starting with it.
	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	sv := doc.ValueAt("A").AsString("illegal")
	assert.Equal(sv, lo.A)
	sv = doc.ValueAt("B/0/A").AsString("illegal")
	assert.Equal(sv, lo.B[0].A)
	sv = doc.ValueAt("/B/1/D/A").AsString("illegal")
	assert.Equal(sv, lo.B[1].D.A)
	sv = doc.ValueAt("/B/2/S").AsString("illegal")
	assert.Equal(sv, "illegal")

	// Now two colons.
	doc, err = gjp.Parse(bs, "::")
	assert.Nil(err)
	sv = doc.ValueAt("A").AsString("illegal")
	assert.Equal(sv, lo.A)
	sv = doc.ValueAt("B::0::A").AsString("illegal")
	assert.Equal(sv, lo.B[0].A)
	sv = doc.ValueAt("B::1::D::A").AsString("illegal")
	assert.Equal(sv, lo.B[1].D.A)

	// Check if is undefined.
	v := doc.ValueAt("you-wont-find-me")
	assert.True(v.IsUndefined())
}

// TestCompare tests comparing two documents.
func TestCompare(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	first, _ := createDocument(assert)
	second := createCompareDocument(assert)
	firstDoc, err := gjp.Parse(first, "/")
	assert.Nil(err)
	secondDoc, err := gjp.Parse(second, "/")
	assert.Nil(err)

	diff, err := gjp.Compare(first, first, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 0)

	diff, err = gjp.Compare(first, second, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 12)
	diff, err = gjp.CompareDocuments(firstDoc, secondDoc, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 12)

	for _, path := range diff.Differences() {
		fv, sv := diff.DifferenceAt(path)
		fvs := fv.AsString("<first undefined>")
		svs := sv.AsString("<second undefined>")
		assert.Different(fvs, svs, path)
	}

	first, err = diff.FirstDocument().MarshalJSON()
	assert.Nil(err)
	second, err = diff.SecondDocument().MarshalJSON()
	assert.Nil(err)
	diff, err = gjp.Compare(first, second, ":")
	assert.Nil(err)
	assert.Length(diff.Differences(), 12)

	// Special case of empty arrays, objects, and null.
	first = []byte(`{}`)
	second = []byte(`{"a":[],"b":{},"c":null}`)

	sdocParsed, err := gjp.Parse(second, "/")
	assert.Nil(err)
	sdocMarshalled, err := sdocParsed.MarshalJSON()
	assert.Nil(err)
	assert.Equal(string(sdocMarshalled), string(second))

	diff, err = gjp.Compare(first, second, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 4)

	first = []byte(`[]`)
	diff, err = gjp.Compare(first, second, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 4)

	first = []byte(`["A", "B", "C"]`)
	diff, err = gjp.Compare(first, second, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 6)

	first = []byte(`"foo"`)
	diff, err = gjp.Compare(first, second, "/")
	assert.Nil(err)
	assert.Length(diff.Differences(), 4)
}

// TestString tests retrieving values as strings.
func TestString(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	sv := doc.ValueAt("A").AsString("illegal")
	assert.Equal(sv, "Level One")
	sv = doc.ValueAt("B/0/B").AsString("illegal")
	assert.Equal(sv, "100")
	sv = doc.ValueAt("B/0/C").AsString("illegal")
	assert.Equal(sv, "true")
	sv = doc.ValueAt("B/0/D/B").AsString("illegal")
	assert.Equal(sv, "10.1")
}

// TestInt tests retrieving values as ints.
func TestInt(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	iv := doc.ValueAt("A").AsInt(-1)
	assert.Equal(iv, -1)
	iv = doc.ValueAt("B/0/B").AsInt(-1)
	assert.Equal(iv, 100)
	iv = doc.ValueAt("B/0/C").AsInt(-1)
	assert.Equal(iv, 1)
	iv = doc.ValueAt("B/0/S/2").AsInt(-1)
	assert.Equal(iv, 1)
	iv = doc.ValueAt("B/0/D/B").AsInt(-1)
	assert.Equal(iv, 10)
}

// TestFloat64 tests retrieving values as float64.
func TestFloat64(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	fv := doc.ValueAt("A").AsFloat64(-1.0)
	assert.Equal(fv, -1.0)
	fv = doc.ValueAt("B/1/B").AsFloat64(-1.0)
	assert.Equal(fv, 200.0)
	fv = doc.ValueAt("B/0/C").AsFloat64(-99)
	assert.Equal(fv, 1.0)
	fv = doc.ValueAt("B/0/S/3").AsFloat64(-1.0)
	assert.Equal(fv, 2.2)
	fv = doc.ValueAt("B/1/D/B").AsFloat64(-1.0)
	assert.Equal(fv, 20.2)
}

// TestBool tests retrieving values as bool.
func TestBool(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	bv := doc.ValueAt("A").AsBool(false)
	assert.Equal(bv, false)
	bv = doc.ValueAt("B/0/C").AsBool(false)
	assert.Equal(bv, true)
	bv = doc.ValueAt("B/0/S/0").AsBool(false)
	assert.Equal(bv, false)
	bv = doc.ValueAt("B/0/S/2").AsBool(false)
	assert.Equal(bv, true)
	bv = doc.ValueAt("B/0/S/4").AsBool(false)
	assert.Equal(bv, true)
}

// TestQuery tests querying a document.
func TestQuery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	bs, _ := createDocument(assert)

	doc, err := gjp.Parse(bs, "/")
	assert.Nil(err)
	pvs, err := doc.Query("Z/*")
	assert.Nil(err)
	assert.Length(pvs, 0)
	pvs, err = doc.Query("*")
	assert.Nil(err)
	assert.Length(pvs, 27)
	pvs, err = doc.Query("/A")
	assert.Nil(err)
	assert.Length(pvs, 1)
	pvs, err = doc.Query("/B/*")
	assert.Nil(err)
	assert.Length(pvs, 24)
	pvs, err = doc.Query("/B/[01]/*")
	assert.Nil(err)
	assert.Length(pvs, 18)
	pvs, err = doc.Query("/B/[01]/*A")
	assert.Nil(err)
	assert.Length(pvs, 4)
	pvs, err = doc.Query("*/S/*")
	assert.Nil(err)
	assert.Length(pvs, 8)
	pvs, err = doc.Query("*/S/3")
	assert.Nil(err)
	assert.Length(pvs, 1)

	pvs, err = doc.Query("/A")
	assert.Nil(err)
	assert.Equal(pvs[0].Path, "/A")
	assert.Equal(pvs[0].Value.AsString(""), "Level One")
}

// TestBuilding tests the creation of documents.
func TestBuilding(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Most simple document.
	doc := gjp.NewDocument("/")
	err := doc.SetValueAt("", "foo")
	assert.Nil(err)

	sv := doc.ValueAt("").AsString("bar")
	assert.Equal(sv, "foo")

	// Positive cases.
	doc = gjp.NewDocument("/")
	err = doc.SetValueAt("/a/b/x", 1)
	assert.Nil(err)
	err = doc.SetValueAt("/a/b/y", true)
	assert.Nil(err)
	err = doc.SetValueAt("/a/c", "quick brown fox")
	assert.Nil(err)
	err = doc.SetValueAt("/a/d/0/z", 47.11)
	assert.Nil(err)
	err = doc.SetValueAt("/a/d/1/z", nil)
	assert.Nil(err)

	iv := doc.ValueAt("a/b/x").AsInt(0)
	assert.Equal(iv, 1)
	bv := doc.ValueAt("a/b/y").AsBool(false)
	assert.Equal(bv, true)
	sv = doc.ValueAt("a/c").AsString("")
	assert.Equal(sv, "quick brown fox")
	fv := doc.ValueAt("a/d/0/z").AsFloat64(8.15)
	assert.Equal(fv, 47.11)
	nv := doc.ValueAt("a/d/1/z").IsUndefined()
	assert.True(nv)

	pvs, err := doc.Query("*x")
	assert.Nil(err)
	assert.Length(pvs, 1)

	// Now provoke errors.
	err = doc.SetValueAt("a", "stupid")
	assert.ErrorMatch(err, ".*corrupt.*")
	err = doc.SetValueAt("a/b/x/y", "stupid")
	assert.ErrorMatch(err, ".*corrupt.*")

	// Legally change values.
	err = doc.SetValueAt("/a/b/x", 2)
	assert.Nil(err)
	iv = doc.ValueAt("a/b/x").AsInt(0)
	assert.Equal(iv, 2)
}

// TestMarshalJSON tests building a JSON document again.
func TestMarshalJSON(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Compare input and output.
	bsIn, _ := createDocument(assert)
	parsedDoc, err := gjp.Parse(bsIn, "/")
	assert.Nil(err)
	bsOut, err := parsedDoc.MarshalJSON()
	assert.Nil(err)
	assert.Equal(bsOut, bsIn)

	// Now create a built one.
	builtDoc := gjp.NewDocument("/")
	err = builtDoc.SetValueAt("/a/2/x", 1)
	assert.Nil(err)
	err = builtDoc.SetValueAt("/a/4/y", true)
	assert.Nil(err)
	bsIn = []byte(`{"a":[null,null,{"x":1},null,{"y":true}]}`)
	bsOut, err = builtDoc.MarshalJSON()
	assert.Nil(err)
	assert.Equal(bsOut, bsIn)
}

//--------------------
// HELPERS
//--------------------

type levelThree struct {
	A string
	B float64
}

type levelTwo struct {
	A string
	B int
	C bool
	D *levelThree
	S []string
}

type levelOne struct {
	A string
	B []*levelTwo
	D time.Duration
	T time.Time
}

func createDocument(assert *asserts.Asserts) ([]byte, *levelOne) {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"1",
					"2.2",
					"true",
				},
			},
			{
				A: "Level Two - 1",
				B: 200,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 20.2,
				},
				S: []string{
					"orange",
					"blue",
					"white",
				},
			},
			{
				A: "Level Two - 2",
				B: 300,
				C: true,
				D: &levelThree{
					A: "Level Three - 2",
					B: 30.3,
				},
			},
		},
		D: 5 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 30, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	assert.Nil(err)
	return bs, lo
}

func createCompareDocument(assert *asserts.Asserts) []byte {
	lo := &levelOne{
		A: "Level One",
		B: []*levelTwo{
			{
				A: "Level Two - 0",
				B: 100,
				C: true,
				D: &levelThree{
					A: "Level Three - 0",
					B: 10.1,
				},
				S: []string{
					"red",
					"green",
					"0",
					"2.2",
					"false",
				},
			},
			{
				A: "Level Two - 1",
				B: 300,
				C: false,
				D: &levelThree{
					A: "Level Three - 1",
					B: 99.9,
				},
				S: []string{
					"orange",
					"blue",
					"white",
					"red",
				},
			},
		},
		D: 10 * time.Second,
		T: time.Date(2018, time.April, 29, 20, 59, 0, 0, time.UTC),
	}
	bs, err := json.Marshal(lo)
	assert.Nil(err)
	return bs
}

// EOF
