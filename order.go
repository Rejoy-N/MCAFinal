// package Order to manage Order data
package main

// import packages
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

// OrderShipAdd Model for Order Shipping data
type OrderShipAdd struct {
		eladdress				string 
		shaddressname 			string	
		shaddressline 			string
		shaddresszip 			string
		shaddresscity			string
		shaddresscountry 		string
}

// Order Model
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

// function SaveOrder is used to save the Order data to the database
func SaveOrderData(o *Order) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists orders (ouid text not null unique, userid text not null, token text not null, orderdetail blob not null, shipdetail blob not null, totalamount text not null, chargedamount integer not null, chargeid text not null, chargestatus text not null, timestamp text, primary key(ouid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into orders (ouid, userid, token, orderdetail, shipdetail, totalamount, chargedamount, chargeid, chargestatus, timestamp) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(o.Ouid, o.Userid, o.Token, o.OrderDetail, o.ShipDetail, o.TotalAmount, o.ChargedAmount, o.Chargeid, o.ChargeStatus, o.Timestamp)
	tx.Commit()
	return err
}

// function Ouid is used to generate the OUID string for a new Order
func Ouid() string {
	id := uuid.NewV4()
	return id.String()
}


// function GetOrders is used to retrieve from the Orders table of the database the Order UUID, Order Value and Order Timestamp 
// for all orders in the entire product table and encode it to JSON
func GetOrders() ([]byte, error) {
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

// function GetOrderFromOuid is used to retrieve from the database the entire order data for the Order UUID table and encode it to JSON
func GetOrderFromOuid(ouid string) (map[string]interface{}, error) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	
	var ou, od, sh string
	var ta, ts int64 
	// var od []byte
	
	q, err := db.Query("select ouid, orderdetail, shipdetail, chargedamount, timestamp from orders where ouid = '" + ouid + "'")
	if err != nil {
		return nil, err
	}
	
	getord := make(map[string]interface{})
	
	for q.Next() {
		q.Scan(&ou, &od, &sh, &ta, &ts)
		getord["ouid"] = ou
		getord["orderdetail"] = od
		getord["shipdetail"] = sh
		getord["chargedamount"] = float64(ta/100)
		getord["timestamp"] = string(time.Unix(ts, 0).Format("02.01.2006 15:04:05"))
	}
	
	var odd ([]map[string]interface{})
	
	err = json.Unmarshal([]byte(od), &odd)
	if err != nil {
		fmt.Println(err)
	}
	
	getord["orderdetail"] = odd
	
	
	osd := make(map[string]interface{})
	
	err = json.Unmarshal([]byte(sh), &osd)
	if err != nil {
		fmt.Println(err)
	}
	
	getord["shipdetail"] = osd
	
	return getord, nil

}




