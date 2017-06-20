// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "log"

var migrations = []string{
	"sql/01-schema.sql",
	"sql/02-schema.sql",
	"sql/03-schema.sql",
}

func (conn *Database) MustRunMigrations(resources *Resources) {

	conn.createMigrationTable()

	applied, err := conn.findAppliedMigrations()
	if err != nil {
		panic(err)
	}

	for _, migration := range migrations {
		run := applied[migration]
		log.Printf("- Migration '%s' has been run? %v\n", migration, run)
		if !run {
			ddl, err := resources.PrivateString(migration)
			if err != nil {
				log.Printf("- Unable to apply: %s.", migration)
				panic(err)
			}
			conn.applyMigration(migration, ddl)
		}
	}

}

func (conn *Database) createMigrationTable() error {
	ddl := `create table if not exists migrations (
						id serial primary key,
						name varchar not null,
						created_at timestamp not null default current_timestamp)`

	_, err := conn.db.Exec(ddl)
	return err
}

func (conn *Database) findAppliedMigrations() (map[string]bool, error) {
	names := make(map[string]bool, 0)

	const q = "select name from migrations"

	rows, err := conn.db.Query(q)
	if err != nil {
		return names, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names[name] = true
	}
	return names, nil
}

func (conn *Database) applyMigration(name, ddl string) {

	tx, err := conn.db.Begin()
	if err != nil {
		log.Printf("- Unable to apply: %s.", name)
		panic(err)
	}

	_, err = conn.db.Exec(ddl)
	if err != nil {
		tx.Rollback()
		log.Printf("- Unable to apply: %s.", name)
		panic(err)
	}

	_, err = conn.db.Exec("insert into migrations (name) values ($1)", name)
	if err != nil {
		tx.Rollback()
		log.Printf("- Unable to apply: %s.", name)
		panic(err)
	}

	tx.Commit()
}
