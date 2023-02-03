package main

import (
	//"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func createMigrationsTable() {
	// table called migrations
	// 3 cols:  ID, name, inserted at
	connStr := "user=postgres password=1234rewQ host=localhost port=5432 dbname=RabbitMQLogger sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
		`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created successfully!")

}

// in  code we need to keep track of all migrations
// migrations: up (new change) : adding cols & down (revert from change) : removing a col I added & key : similar to commit message

/*
Migration will be in code


the record of your migration will be on psql


==================Flow====================
on start up

1- Query migration table from postgres
2- Compare with the migrations in code
3- if there is a mismatch
4- update and apply the migrations
*/

type Migration struct {
	Key  string // commit message
	Up   string // sql Query
	Down string // sql Query
}

var Migrations = []Migration{
	{ //
		Key: "initial Schema",
		Up: `
			CREATE TABLE logs(
				id SERIAL PRIMARY KEY,
				body text,
				created_at TIMESTAMPZ DEFAULT now(),
				updated_at TIMESTAMPZ DEDAULT now()
			);
		`,
		Down: `
			DROP TABLE logs;
		`,
	},
}

func runMigrations(
	dbctx sql.DB,
) error {

	// rows, err := dbctx.QueryContext(
	// 	context.TODO(),
	// 	``,
	// )

	// TODO: Handle error

	//map rows to []string, hint: ur selecting name
	keys := []string{}
	// for rows.Next() {
	// 	var key string
	// 	rows.Scan(&key) // MUST pass reference to the var
	// }

	if len(keys) == len(Migrations) {
		return nil
	}

	for i := 0; i < len(Migrations); i++ {
		// compare with the key array
		// if it matches ; continue
		// if it doesnt do the up of the migration
	}

	return nil
}
