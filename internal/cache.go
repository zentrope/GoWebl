// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"net"
	"sync"
)

// alternative: https://github.com/orcaman/concurrent-map
var DNS_CACHE = make(map[string]string, 0)
var DNS_MUTEX = sync.Mutex{}

// TODO: thread safe
// TODO: go thread to age out values
func DNSLookup(address string) string {

	if address == "0.0.0.0" {
		return address
	}

	DNS_MUTEX.Lock()
	defer DNS_MUTEX.Unlock()

	if name, found := DNS_CACHE[address]; found {
		return name
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

	DNS_CACHE[address] = names[0]
	return names[0]
}
