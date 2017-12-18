package migrations

import (
	"fmt"
	"github.com/go-pg/migrations"
	"project/models"
	"time"
)

func init() {
	migrations.Register(
		func(db migrations.DB) error {
			fmt.Println("Seeding users...")
			user1 := models.User{
				Username:  "Alice",
				Email:     "alice@example.com",
				CreatedAt: time.Now(),
			}
			user2 := models.User{
				Username:  "Bob",
				Email:     "bob@example.com",
				CreatedAt: time.Now(),
			}
			err := db.Insert(&user1, &user2)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("Truncating users table...")
			_, err := db.Exec(`TRUNCATE users`)
			return err
		})
}
