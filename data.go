//package User for managing User data
package main

// import packages
import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User Model 
type User struct {
	Uuid     string            `valid:"required,uuidv4"`
	Username string            `valid:"required,alphanum"`
	Password string            `valid:"required"`
	Fname    string            `valid:"required,alpha"`
	Lname    string            `valid:"required,alpha"`
	Email    string            `valid:"required,email"`
	Errors   map[string]string `valid:"-"`
}

// function SaveData saves the accepted user enetered data to the database
func SaveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (uuid text not null unique, firstname text not null, lastname text not null, username text not null unique, email text not null unique, password text not null, primary key(uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Uuid, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

// function UserExists is used to check if a user provided username exists in the database and its subsequent authentication
// If the username exists in the database, the user provided password is compared with the password saved in the database
func UserExists(u *User) (bool, string) {
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

// function CheckUser is used to check if a user provided username exists in the database
func CheckUser(user string) bool {
	if user =="" {
		return false
	}
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var un string
	fmt.Println(user)
	q, err := db.Query("select username from users where username = '" + user + "'")
	if err != nil {
		return false
	}
	{
	for q.Next() {
		q.Scan(&un)
		}
	}	
	if un == user {
		return true
	}
	return false
}

// function GetUserFromUuid is used to retrieve from the database user information corresponding to the UUID  
func GetUserFromUuid(uuid string) *User {
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

// function getUserEmailFromUuid is the used to retrieve from the database the email address corresponding to the user's UUID
func GetUserEmailFromUuid(uuid string) string {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var em string
	q, err := db.Query("select email from users where uuid = '" + uuid + "'")
	if err != nil {
		return em
	}
	for q.Next() {
		q.Scan(&em)
	}
	return em
}

// function EncryptPass is used to encrypt user password 
func EncryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

// function Uuid is used to generate the UUID string for a new user
func Uuid() string {
	id := uuid.NewV4()
	return id.String()
}
