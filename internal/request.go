//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package internal

import (
	"log"
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

		r.Host = DNSLookup(r.Address)
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
