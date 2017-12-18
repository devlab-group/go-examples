package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/go-pg/pg"
	"os"
	"project/migrations"
	"project/models"
	"time"
)

var db *pg.DB

func main() {
	db = pg.Connect(dbConnectOptions())
	defer db.Close()

	switch os.Args[1] {
	case "users":
		resolveUserCommands()
	case "migrations":
		migrations.Run(db)
	default:
		fatal("Unknown command")
	}
}

func resolveUserCommands() {
	switch os.Args[2] {
	case "add":
		if len(os.Args) != 5 {
			fatal("Usage: users add USERNAME EMAIL")
		}
		user, err := createUser(db, os.Args[3], os.Args[4])
		checkErr(err)
		fmt.Printf("New user with id %d was created successfully\n", user.Id)
	case "del":
		if len(os.Args) < 4 {
			fatal("Usage: users del ID...")
		}
		deletedUsers, err := removeUser(db, os.Args[3:])
		checkErr(err)
		fmt.Printf("%d users were deleted successfully\n", deletedUsers)
	case "update":
		if len(os.Args) > 5 {
			fatal("Usage: users update ID EMAIL USERNAME")
		}
		user, err := updateUser(db, os.Args[3], os.Args[4:])
		checkErr(err)
		fmt.Println(user)
	case "all":
		if len(os.Args) > 3 {
			fatal("Usage: users all")
		}
		users, err := allUsers(db)
		checkErr(err)

		fmt.Println("id | username | email | created_at")

		for _, user := range users {
			fmt.Printf("%3v | %8v | %15v | %20v\n", user.Id, user.Username, user.Email, user.CreatedAt)
		}
	}
}

/**
 * Creates new user with passed username and email
 */
func createUser(db *pg.DB, username, email string) (user models.User, err error) {
	user = models.User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
	err = db.Insert(&user)
	return
}

/**
 * Removes users with by ids
 */
func removeUser(db *pg.DB, ids []string) (int, error) {
	inIds := pg.In(ids)
	res, err := db.Model(&models.User{}).Where("id IN (?)", inIds).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

/**
 * Updates user by id
 */
func updateUser(db *pg.DB, id string, fields []string) (models.User, error) {
	var user models.User
	model := db.Model(&user)

	if len(fields) > 0 {
		model.Set("email = ?", fields[0])
	}

	if len(fields) > 1 {
		model.Set("username = ?", fields[1])
	}
	_, err := model.
		Where("id = ?", id).
		Returning("*").
		Update()

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

/**
 * Returns all users
 */
func allUsers(db *pg.DB) ([]models.User, error) {
	var users []models.User
	err := db.Model(&users).Select()
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

/**
 * Parse database config .ini file and return connect options
 */
func dbConnectOptions() *pg.Options {
	cfg, err := ini.Load("db.ini")
	checkErr(err)

	section, err := cfg.GetSection("postgres")
	checkErr(err)

	dbname, err := section.GetKey("dbname")
	checkErr(err)

	user, err := section.GetKey("user")
	checkErr(err)

	password, err := section.GetKey("password")
	checkErr(err)

	addr := section.Key("addr").Validate(func(in string) string {
		if len(in) == 0 {
			return "localhost:5432"
		}
		return in
	})

	return &pg.Options{
		User:     user.String(),
		Password: password.String(),
		Database: dbname.String(),
		Addr:     addr,
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}
