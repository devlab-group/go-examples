package models

import (
	"fmt"
	"time"
)

/**
 * User model
 */
type User struct {
	Id        int64
	Username  string    `sql:",notnull,unique"`
	Email     string    `sql:",notnull,unique"`
	CreatedAt time.Time `sql:",notnull"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %s %s>", u.Id, u.Username, u.Email, u.CreatedAt)
}
