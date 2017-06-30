// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"log"
	"net"
	"net/http"
	"time"
)

type RequestStat struct {
	Address      string
	Host         string
	DateRecorded time.Time
	Method       string
	Path         string
	UserAgent    string
	Referer      string
}

func (conn *Database) RecordRequest(r *http.Request) {

	address := r.Header.Get("X-Forwarded-For")

	if address == "" {
		address = "0.0.0.0"
	}

	stat := &RequestStat{
		Address:      address,
		DateRecorded: time.Now(),
		Method:       r.Method,
		Path:         r.RequestURI,
		UserAgent:    r.UserAgent(),
		Referer:      r.Referer(),
	}

	// Revisit this if the site actually gets a load.
	go func() {
		if err := conn.writeRequest(stat); err != nil {
			log.Printf("ERROR: recording stat err: %v, stat: %#v", err, stat)
		}
	}()
}

func findHost(address string, cache map[string]string) string {

	if name, found := cache[address]; found {
		return name
	}

	if address == "0.0.0.0" {
		return address
	}

	names, err := net.LookupAddr(address)
	if err != nil {
		return address
	}

	if len(names) == 0 {
		cache[address] = address
		return address
	}

	name := names[0]

	cache[address] = name

	return name
}

func (conn *Database) RecentRequests(limit int) ([]*RequestStat, error) {
	q := `select address, date_recorded, method, path, user_agent, referer from
					request where abbrev(address) <> '0.0.0.0'
						order by date_recorded desc limit $1`

	rows, err := conn.db.Query(q, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := make([]*RequestStat, 0)

	tempCache := make(map[string]string, 0)

	for rows.Next() {
		var r RequestStat
		err := rows.Scan(
			&r.Address,
			&r.DateRecorded,
			&r.Method,
			&r.Path,
			&r.UserAgent,
			&r.Referer,
		)

		if err != nil {
			return nil, err
		}

		r.Host = findHost(r.Address, tempCache)
		// names, err := net.LookupAddr(r.Address)
		// if err != nil {
		//	r.Host = r.Address
		// } else {
		//	r.Host = fmt.Sprintf("%v", names)
		// }

		requests = append(requests, &r)
	}

	return requests, nil
}

func (conn *Database) writeRequest(r *RequestStat) error {

	q := `insert into request
					(address, date_recorded, method, path, user_agent, referer)
						values ($1, $2, $3, $4, $5, $6)`

	_, err := conn.db.Exec(q,
		r.Address,
		r.DateRecorded,
		r.Method,
		r.Path,
		r.UserAgent,
		r.Referer,
	)

	return err
}
