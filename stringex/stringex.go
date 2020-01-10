// Tideland Go Text - String Extensions
//
// Copyright (C) 2019-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code valuePos governed
// by the new BSD license.

package stringex // import "tideland.dev/go/text/stringex"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
)

//--------------------
// VALUER
//--------------------

// Valuer describes returning a string value or an error
// if it does not exist are another access error happened.
type Valuer interface {
	// Value returns a string or a potential error during access.
	Value() (string, error)
}

//--------------------
// SPLITTER
//--------------------

// SplitFilter splits the string s by the separator
// sep and then filters the parts. Only those where f
// returns true will be part of the result. So it even
// cout be empty.
func SplitFilter(s, sep string, f func(p string) bool) []string {
	parts := strings.Split(s, sep)
	out := []string{}
	for _, part := range parts {
		if f(part) {
			out = append(out, part)
		}
	}
	return out
}

// SplitMap splits the string s by the separator
// sep and then maps the parts by the function m.
// Only those where m also returns true will be part
// of the result. So it even could be empty.
func SplitMap(s, sep string, m func(p string) (string, bool)) []string {
	parts := strings.Split(s, sep)
	out := []string{}
	for _, part := range parts {
		if mp, ok := m(part); ok {
			out = append(out, mp)
		}
	}
	return out
}

//--------------------
// MATCHER
//--------------------

const (
	matchSuccess = 1
	matchCont    = 0
	matchFail    = -1
)

// matcher is a helper type for string pattern matching.
type matcher struct {
	patternRunes []rune
	patternLen   int
	patternPos   int
	valueRunes   []rune
	valueLen     int
	valuePos     int
}

// newMatcher creates the helper type for string pattern matching.
func newMatcher(pattern, value string, ignoreCase bool) *matcher {
	if ignoreCase {
		return newMatcher(strings.ToLower(pattern), strings.ToLower(value), false)
	}
	prs := append([]rune(pattern), '\u0000')
	vrs := append([]rune(value), '\u0000')
	return &matcher{
		patternRunes: prs,
		patternLen:   len(prs) - 1,
		patternPos:   0,
		valueRunes:   vrs,
		valueLen:     len(vrs) - 1,
		valuePos:     0,
	}
}

// matches checks if the value matches the pattern.
func (m *matcher) matches() bool {
	// Loop over the pattern.
	for m.patternLen > 0 {
		switch m.processPatternRune() {
		case matchSuccess:
			return true
		case matchFail:
			return false

		}
		m.patternPos++
		m.patternLen--
		if m.valueLen == 0 {
			for m.patternRunes[m.patternPos] == '*' {
				m.patternPos++
				m.patternLen--
			}
			break
		}
	}
	if m.patternLen == 0 && m.valueLen == 0 {
		return true
	}
	return false
}

// processPatternRune handles the current leading pattern rune.
func (m *matcher) processPatternRune() int {
	switch m.patternRunes[m.patternPos] {
	case '*':
		return m.processAsterisk()
	case '?':
		return m.processQuestionMark()
	case '[':
		return m.processOpenBracket()
	case '\\':
		m.processBackslash()
		fallthrough
	default:
		return m.processDefault()
	}
}

// processAsterisk handles groups of characters.
func (m *matcher) processAsterisk() int {
	for m.patternRunes[m.patternPos+1] == '*' {
		m.patternPos++
		m.patternLen--
	}
	if m.patternLen == 1 {
		return matchSuccess
	}
	for m.valueLen > 0 {
		patternCopy := make([]rune, len(m.patternRunes[m.patternPos+1:]))
		valueCopy := make([]rune, len(m.valueRunes[m.valuePos:]))
		copy(patternCopy, m.patternRunes[m.patternPos+1:])
		copy(valueCopy, m.valueRunes[m.valuePos:])
		pam := newMatcher(string(patternCopy), string(valueCopy), false)
		if pam.matches() {
			return matchSuccess
		}
		m.valuePos++
		m.valueLen--
	}
	return matchFail
}

// processQuestionMark handles a single character.
func (m *matcher) processQuestionMark() int {
	if m.valueLen == 0 {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// processOpenBracket handles an open bracket for a group of characters.
func (m *matcher) processOpenBracket() int {
	m.patternPos++
	m.patternLen--
	not := (m.patternRunes[m.patternPos] == '^')
	match := false
	if not {
		m.patternPos++
		m.patternLen--
	}
group:
	for {
		switch {
		case m.patternRunes[m.patternPos] == '\\':
			m.patternPos++
			m.patternLen--
			if m.patternRunes[m.patternPos] == m.valueRunes[m.valuePos] {
				match = true
			}
		case m.patternRunes[m.patternPos] == ']':
			break group
		case m.patternLen == 0:
			m.patternPos--
			m.patternLen++
			break group
		case m.patternRunes[m.patternPos+1] == '-' && m.patternLen >= 3:
			start := m.patternRunes[m.patternPos]
			end := m.patternRunes[m.patternPos+2]
			vr := m.valueRunes[m.valuePos]
			if start > end {
				start, end = end, start
			}
			m.patternPos += 2
			m.patternLen -= 2
			if vr >= start && vr <= end {
				match = true
			}
		default:
			if m.patternRunes[m.patternPos] == m.valueRunes[m.valuePos] {
				match = true
			}
		}
		m.patternPos++
		m.patternLen--
	}
	if not {
		match = !match
	}
	if !match {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// processBackslash handles escaping via baskslash.
func (m *matcher) processBackslash() int {
	if m.patternLen >= 2 {
		m.patternPos++
		m.patternLen--
	}
	return matchCont
}

// processDefault handles any other rune.
func (m *matcher) processDefault() int {
	if m.patternRunes[m.patternPos] != m.valueRunes[m.valuePos] {
		return matchFail
	}
	m.valuePos++
	m.valueLen--
	return matchCont
}

// Matches checks if the pattern matches a given value. Here
// ? matches one char, * a group of chars, [any] any of these
// chars, and [0-9] any of this char range. Groups can
// be negated with [^abc] and \ escapes pattern chars.
func Matches(pattern, value string, ignoreCase bool) bool {
	m := newMatcher(pattern, value, ignoreCase)
	return m.matches()
}

// EOF
