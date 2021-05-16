// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package dj

// nodeLen returns the length of the passed node (which can be a single
// value too).
func nodeLen(node interface{}) int {
	if node == nil {
		return 0
	}
	switch n := node.(type) {
	case []interface{}:
		return len(n)
	case map[string]interface{}:
		return len(n)
	case string, int, float64, bool:
		return 1
	}
	return 0
}

// EOF
