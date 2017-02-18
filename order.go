package main

import (
	"database/sql"
	"encoding/json"
	"time"
	// "strconv"
	// "reflect"
	"fmt"
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
	ChargedAmount   uint64						`valid:"-"` 
	Chargeid	    string						`valid:"-"`
	ChargeStatus	string						`valid:"-"`
	Timestamp		string						`valid:"-"`
	Errors   		map[string]string 			`valid:"-"`
}

func saveOrderData(o *Order) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists orders (ouid text not null unique, userid text not null, token text not null, orderdetail blob not null, shipdetail blob not null, totalamount text not null, chargedamount integer not null, chargeid text not null, chargestatus text not null, timestamp text, primary key(ouid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into orders (ouid, userid, token, orderdetail, shipdetail, totalamount, chargedamount, chargeid, chargestatus, timestamp) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(o.Ouid, o.Userid, o.Token, o.OrderDetail, o.ShipDetail, o.TotalAmount, o.ChargedAmount, o.Chargeid, o.ChargeStatus, o.Timestamp)
	tx.Commit()
	return err
}

func Ouid() string {
	id := uuid.NewV4()
	return id.String()
}

func getOrders() ([]byte, error) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var ou string
	var ta, ts int64 
	q, err := db.Query("select ouid, chargedamount, timestamp from orders")
	if err != nil {
		fmt.Println(err)
	}
	
	var a []interface{}
	
	for q.Next() {
		q.Scan(&ou, &ta, &ts)
		b := make(map[string]interface{})	
		b["ouid"] = ou
		b["chargedamount"] = float64(ta)/100
		// b["timestamp"] = ts
		b["timestamp"] = string(time.Unix(ts, 0).Format("02.01.2006 15:04:05"))
		a = append(a, b)
	}
	
	getord, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(getord))
		return getord, nil
}

/*
func getOrderFromOuid(ouid string) ([]byte, error) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var ou, od, sh, ta, ts string
	q, err := db.Query("select ouid, orderdetail, shipdetail, chargedamount, timestamp from orders where ouid = '" + ouid + "'")
	if err != nil {
		return err
	}
	for q.Next() {
		q.Scan(&ou, &od, &sh, &ta, &ts)
		getord := make(map[string]interface{})
		getord["ouid"] = ou
		getord["orderdetail"] = od
		getord["shipdetail"] = sh
		getord["chargedamount"] = ta
		getord["timestamp"] = ts
	}
	
	jsdata, err := json.Marshal(getord)
	if err != nil {
		return nil, err
	}
	//  fmt.Println(string(jsData))
		return jsData, nil	
}

*/
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

