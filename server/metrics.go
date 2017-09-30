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

package server

import (
	"database/sql"
)

type Metric struct {
	Key   string
	Value int
}

const topHits = `
	select address as key, count(*) as value
		from request where address <> '0.0.0.0'
			group by address order by value desc`

const topRoutes = `
	select method || ' ' || substring(path, 0, 100) as key, count(*) as value
		from request group by key order by value desc`

const topRefers = `
	select substring(referer, 0, 100) as key, count(*) as value
		from request where referer <> '' group by key order by value desc`

const hitsPerDay = `
	select left(date_trunc('day', date_recorded)::text, 10) as key, count (*) as value
		from request group by key order by key desc`

func (conn *Database) HitsPerDay() ([]*Metric, error) {
	return conn.runMetric(hitsPerDay)
}

func (conn *Database) TopHits() ([]*Metric, error) {
	metrics, err := conn.runMetric(topHits)
	if err != nil {
		return nil, err
	}

	for i, m := range metrics {
		name := DNSLookup(m.Key)
		metrics[i].Key = name
	}

	return metrics, nil
}

func (conn *Database) TopRoutes() ([]*Metric, error) {
	return conn.runMetric(topRoutes)
}

func (conn *Database) TopRefers() ([]*Metric, error) {
	return conn.runMetric(topRefers)
}

func (conn *Database) runMetric(query string) ([]*Metric, error) {
	rows, err := conn.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	metrics := make([]*Metric, 0)

	for rows.Next() {
		metric, err := rowToMetric(rows)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func rowToMetric(rows *sql.Rows) (*Metric, error) {
	var m Metric
	err := rows.Scan(&m.Key, &m.Value)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
