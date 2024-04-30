package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	migrationsPath := filepath.Join("db", "postgres", "migration", "*.up.sql")
	upMigrations, err := filepath.Glob(migrationsPath)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(migrationsPath)
	fmt.Println(upMigrations)
}
