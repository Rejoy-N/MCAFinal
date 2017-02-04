package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid     string            `valid:"required,uuidv4"`
	Username string            `valid:"required,alphanum"`
	Password string            `valid:"required"`
	Fname    string            `valid:"required,alpha"`
	Lname    string            `valid:"required,alpha"`
	Email    string            `valid:"required,email"`
	Errors   map[string]string `valid:"-"`
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (uuid text not null unique, firstname text not null, lastname text not null, username text not null unique, email text not null unique, password text not null, primary key(uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Uuid, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

func userExists(u *User) (bool, string) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var ps, uu string
	q, err := db.Query("select uuid, password from users where username = '" + u.Username + "'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&uu, &ps)
	}
	pw := bcrypt.CompareHashAndPassword([]byte(ps), []byte(u.Password))
	if uu != "" && pw == nil {
		return true, uu
	}
	return false, ""
}

func checkUser(user string) bool {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var un string
	q, err := db.Query("select username from users where username = '" + user + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&un)
	}
	if un == user {
		return true
	}
	return false
}

func getUserFromUuid(uuid string) *User {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var uu, fn, ln, un, em, pass string
	q, err := db.Query("select * from users where uuid = '" + uuid + "'")
	if err != nil {
		return &User{}
	}
	for q.Next() {
		q.Scan(&uu, &fn, &ln, &un, &em, &pass)
	}
	return &User{Username: un, Fname: fn, Lname: ln, Email: em, Password: pass}
}

func encryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

func Uuid() string {
	id := uuid.NewV4()
	return id.String()
}
