package internal

import "log"

//
// createuser blogsvc --createdb --superuser -P
// createdb blogdb -O blogsvc
//

var migrations = []string{
	"sql/01-schema.sql",
}

func (conn *Database) MustRunMigrations(resources *Resources) {

	conn.createMigrationTable()

	applied, err := conn.findAppliedMigrations()

	for _, migration := range migrations {
		run := applied[migration]
		log.Printf("- Migration '%s' has been run? %v\n", migration, run)
		if !run {
			ddl, err := resources.PrivateString(migration)
			if err != nil {
				panic(err)
			}
			conn.applyMigration(migration, ddl)
		}
	}

	if err != nil {
		panic(err)
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
		panic(err)
	}

	_, err = conn.db.Exec(ddl)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = conn.db.Exec("insert into migrations (name) values ($1)", name)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
