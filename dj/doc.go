// Tideland Go Text - Dynamic JSON
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package dj provides the dynmic work with JSON documents. It can
// parse and create documents and access or add values inside of it. Also
// removing is possible.
//
//     myCustomer, err := dj.Parse(aCustomerReader)
//     if err != nil {
//         ...
//     }
//     firstStreet := myCustomer.At("addresses", "#0", "street").AsString()
//
// The value passed to AsString() will panic if an access does not match (the
// hard way) or return the default value for the type if the value is nil. And
// there are methods to set values.
package dj // import "tideland.dev/go/text/dj"

// EOF
