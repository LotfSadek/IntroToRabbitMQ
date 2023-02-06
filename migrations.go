package main

import (
	//"context"
	"database/sql"
	"fmt"

	// "log"
	// "time"

	_ "github.com/lib/pq"
)

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
	Key  string
	Up   string
	Down string
}

var Migrations = []Migration{
	{
		Key: "create_logs_table",
		Up: `
			CREATE TABLE IF NOT EXISTS logs (
				id serial PRIMARY KEY,
				body text NOT NULL,
				created_at timestamptz NOT NULL DEFAULT NOW(),
				inserted_at timestamptz NOT NULL DEFAULT NOW()
			)
		`,
		Down: `
			DROP TABLE logs
		`,
	},
	{
		Key: "testingFirstField",
		Up: `
			INSERT into logs(body) values ('test');
		`,
		Down: `
			DROP TABLE logs
		`,
	},
}

func migrationsManager() {
	// loading the config file
	config, err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// Create the migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			key text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Check the number of migrations in the database
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM migrations
	`).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Number of rows in migrations table in DB:", count)

	// If there's a mismatch between the number of migrations in the database and the array of migrations, apply the missing migrations
	if count != len(Migrations) {
		fmt.Printf("Mismatch detected: %d migrations in database, %d migrations in code\n", count, len(Migrations))

		// Apply the missing migrations
		for i := count; i < len(Migrations); i++ {
			fmt.Printf("Applying migration %d: %s\n", i+1, Migrations[i].Key)
			_, err = db.Exec(Migrations[i].Up)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = db.Exec(`
				INSERT INTO migrations (key)
				VALUES ($1)
			`, Migrations[i].Key)
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}

}

// func runMigrations(
// 	//dbctx sql.DB,
// ) error {

// 	// rows, err := dbctx.QueryContext(
// 	// 	context.TODO(),
// 	// 	``,
// 	// )

// 	// TODO: Handle error

// 	//map rows to []string, hint: ur selecting name
// 	keys := []string{}
// 	// for rows.Next() {
// 	// 	var key string
// 	// 	rows.Scan(&key) // MUST pass reference to the var
// 	// }

// 	if len(keys) == len(Migrations) {
// 		return nil
// 	}

// 	for i := 0; i < len(Migrations); i++ {
// 		// compare with the key array
// 		// if it matches ; continue
// 		// if it doesnt do the up of the migration
// 	}

// 	return nil
// }
