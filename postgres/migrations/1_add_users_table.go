package migrations

import (
	"fmt"
	"github.com/go-pg/migrations"
	"project/models"
)

func init() {
	migrations.Register(
		func(db migrations.DB) error {
			fmt.Println("Creating users table...")
			err := db.CreateTable(&models.User{}, nil)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("Dropping table users...")
			_, err := db.Exec(`DROP TABLE users`)
			return err
		})
}
