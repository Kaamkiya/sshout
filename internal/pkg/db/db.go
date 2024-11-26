package db

import (
	"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	Username string
	PasswordHash string
	Joined time.Time
}

func Open() (err error) {
	db, err = sql.Open("sqlite3", "sshout.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username      TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		joined        TEXT NOT NULL
	);`)
	
	return err
}

func GetUser(username string) (User, error) {
	var passwordHash string
	var joinedString string

	row := db.QueryRow(`SELECT password_hash, joined FROM users WHERE username = ?`, username)
	err := row.Scan(&passwordHash, &joinedString)

	joinedTime, _ := time.Parse(time.UnixDate, joinedString)

	u := User{
		Username: username,
		PasswordHash: passwordHash,
		Joined: joinedTime,
	}

	return u, err
}

func AddUser(username, passwordHash string) (User, error) {
	user := User{
		Username: username,
		PasswordHash: passwordHash,
		Joined: time.Now(),
	}

	_, err := db.Exec(
		`INSERT INTO users(username, password_hash, joined) VALUES (?, ?, ?)`,
		user.Username,
		user.PasswordHash,
		user.Joined.Format(time.UnixDate),
	)
	return user, err
}

func Close() error {
	return db.Close()
}
