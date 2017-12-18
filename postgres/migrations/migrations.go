package migrations

import (
	"flag"
	"fmt"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"os"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
Usage:
  go run *.go <command> [args]
`

func Run(db *pg.DB) {
	flag.Usage = usage
	flag.Parse()

	// Take arguments starting from second due to this call goes through main package
	oldVersion, newVersion, err := migrations.Run(db, flag.Args()[2:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Printf(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}
