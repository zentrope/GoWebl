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
	"fmt"

	// DB drivers require a silent import.
	_ "github.com/lib/pq"
)

// Database represents a connection to a database.
type Database struct {
	Config StorageConfig
	db     *sql.DB
}

// NewDatabase returns a connection use can use to make queries.
func NewDatabase(config StorageConfig) *Database {
	return &Database{config, nil}
}

// MustConnect starts a connection to the database, or panics if it
// doesn't succeed.
func (conn *Database) MustConnect() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conn.Config.User, conn.Config.Password, conn.Config.Database,
		conn.Config.Host, conn.Config.Port)
	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}

	conn.db = db
}

// Disconnect terminates the connection to the database.
func (conn *Database) Disconnect() {
	conn.db.Close()
}
