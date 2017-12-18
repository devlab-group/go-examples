package main

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type User struct {
	id                          int
	username, phone, created_at string
}

func init() {
	db, err := sql.Open("postgres", dbConnectParams())
	checkErr(err)
	defer db.Close()

	/**
	 * Create `users` table on init
	 */
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
	     id serial NOT NULL,
	     username character varying(100) NOT NULL,
       phone character varying(100),
	     created_at date,
	     CONSTRAINT user_pkey PRIMARY KEY (id)
	   )
	   WITH (OIDS=FALSE)`)
	checkErr(err)
}

func main() {
	db, err := sql.Open("postgres", dbConnectParams())
	checkErr(err)
	defer db.Close()

	switch os.Args[2] {
	case "add":
		if len(os.Args) != 5 {
			fatal("Usage: users add USERNAME PHONE")
		}
		id, err := insert(db, os.Args[3], os.Args[4])
		checkErr(err)
		fmt.Printf("New user with id %d was created successfully\n", id)
	case "del":
		if len(os.Args) < 4 {
			fatal("Usage: users del ID...")
		}
		err := remove(db, os.Args[3:])
		checkErr(err)
		fmt.Println("Users were deleted successfully\n")
	case "update-phone":
		if len(os.Args) != 5 {
			fatal("Usage: users update-phone ID PHONE")
		}
		err := updatePhone(db, os.Args[3], os.Args[4])
		checkErr(err)
		fmt.Println("Phone was updated successfully")
	case "show":
		if len(os.Args) > 4 {
			fatal("Usage: users show [SUBSTRING]")
		}
		var s string
		if len(os.Args) == 4 {
			s = os.Args[3]
		}
		res, err := show(db, s)
		checkErr(err)

		fmt.Println("id | username | phone | created_at")

		for _, user := range res {
			fmt.Printf("%3v | %8v | %15v | %20v\n", user.id, user.username, user.phone, user.created_at)
		}
	}
}

/**
 * Creates new user with passed username and phone
 * Raw Query
 */
func insert(db *sql.DB, username, phone string) (insertedId int, err error) {
	err = db.QueryRow("INSERT INTO users(username,phone,created_at) VALUES ($1,$2,$3) returning id;", "Molly", "+444567564543", time.Now()).Scan(&insertedId)
	return
}

/**
 * Removes user with passed ids
 * Statement
 */
func remove(db *sql.DB, ids []string) error {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, v := range ids {
		_, err = stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * Updates phone number of user with passed id
 * Transaction
 */
func updatePhone(db *sql.DB, id, phone string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("UPDATE users SET phone = $1 WHERE id = $2;", phone, id)

	if err != nil {
		return err
	}
	return tx.Commit()
}

func show(db *sql.DB, arg string) ([]User, error) {
	var s string
	if len(arg) != 0 {
		s = "WHERE username LIKE '%" + arg + "%'"
	}
	rows, err := db.Query("SELECT * FROM users " + s + " ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]User, 0)
	var user User
	for rows.Next() {
		err = rows.Scan(&user.id, &user.username, &user.phone, &user.created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

/**
 * Parse database config .ini file and return formatted connection string
 */
func dbConnectParams() string {
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

	host := section.Key("host").Validate(func(in string) string {
		if len(in) == 0 {
			return "localhost"
		}
		return in
	})
	port := section.Key("port").Validate(func(in string) string {
		if len(in) == 0 {
			return "5432"
		}
		return in
	})

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)
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
