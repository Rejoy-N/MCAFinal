package main

import (
	"database/sql"
	// "encoding/json"
	// "reflect"
	// "fmt"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	// "golang.org/x/crypto/bcrypt"
)

type OrderShipAdd struct {
		eladdress				string 
		shaddressname 			string	
		shaddressline 			string
		shaddresszip 			string
		// shaddressstate 
		shaddresscity			string
		shaddresscountry 		string
}

type Order struct {
	Ouid     		string            			`valid:"required,uuidv4"`
	Userid    		string           			`valid:"required,alphanum"`
	Token 			string						`valid:"required,alphanum"`
	OrderDetail     string 						`valid:"-"`
	ShipDetail		string						`valid:"-"`
	TotalAmount     string						`valid:"-"`
	Chargeid	    string						`valid:"-"`
	ChargeStatus	string						`valid:"-"`
	Timestamp		string						`valid:"-"`
	Errors   		map[string]string 			`valid:"-"`
}

func saveOrderData(o *Order) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists orders (ouid text not null unique, userid text not null, token text not null, orderdetail blob not null, shipdetail blob not null, totalamount text not null, chargeid text not null, chargestatus text not null, timestamp text, primary key(ouid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into orders (ouid, userid, token, orderdetail, shipdetail, totalamount, chargeid, chargestatus, timestamp) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(o.Ouid, o.Userid, o.Token, o.OrderDetail, o.ShipDetail, o.TotalAmount, o.Chargeid, o.ChargeStatus, o.Timestamp)
	tx.Commit()
	return err
}

func Ouid() string {
	id := uuid.NewV4()
	return id.String()
}

/*
		type Ordata struct {
			Puid     string `json:"puid"`
			Pname    string `json:"pname"`
			Quantity string `json:"quantity"`
			Price    string `json:"price"`
			Image    string `json:"image"`
		}
		
		type Orderdata struct {
			Token string
			Userid string
			Ord []Ordata
		}		

		var Or []Ordata		
		
		err = json.Unmarshal(Orderdata, &Or)
		if err != nil {
			fmt.Println(err)
		}
		
		var Od Orderdata
		
		Od.Token := token
		Od.Userid := uuid 		
		Od.Ord := Or
		*/

