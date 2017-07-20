// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "net"

var DNS_CACHE = make(map[string]string, 0)

// TODO: thread safe
// TODO: go thread to age out values
func DNSLookup(address string) string {

	if name, found := DNS_CACHE[address]; found {
		return name
	}

	if address == "0.0.0.0" {
		return address
	}

	names, err := net.LookupAddr(address)
	if err != nil {
		DNS_CACHE[address] = address
		return address
	}

	if len(names) == 0 {
		DNS_CACHE[address] = address
		return address
	}

	name := names[0]

	DNS_CACHE[address] = name

	return name
}
