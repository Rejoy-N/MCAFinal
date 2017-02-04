package main

import (
	"database/sql"
	"encoding/json"
	// "reflect"
	// "fmt"
	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	// "golang.org/x/crypto/bcrypt"
)

type Product struct {
	Puid     string            `valid:"required,uuidv4"`
	Pname    string            `valid:"required,alphanum"`
	Quantity string            `valid:"required,numeric"`
	Price    string            `valid:"required,numeric"`
	Image    string            `valid:"-"`
	Errors   map[string]string `valid:"-"`
}

func saveProductData(p *Product) error {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists products (puid text not null unique, pname text not null, quantity text not null, price text not null , image text, primary key(puid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into products (puid, pname, quantity, price, image ) values (?, ?, ?, ? , ? )")
	_, err := stmt.Exec(p.Puid, p.Pname, p.Quantity, p.Price, p.Image)
	tx.Commit()
	return err
}

func productExists(p *Product) (bool, string) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var pid, pn string
	q, err := db.Query("select puid, pname from products where pname = '" + p.Pname + "'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&pid, &pn)
	}
	if pid != "" && pn != "" {
		return true, pid
	}
	return false, ""
}

func checkProduct(product string) bool {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var pr string
	q, err := db.Query("select product from products where pname = '" + product + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&pr)
	}
	if pr == product {
		return true
	}
	return false
}

func getProductFromPuid(puid string) *Product {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	var pid, product, quantity, price, image string
	q, err := db.Query("select * from products where puid = '" + puid + "'")
	if err != nil {
		return &Product{}
	}
	for q.Next() {
		q.Scan(&pid, &product, &quantity, &price, &image)
	}
	return &Product{Pname: product, Quantity: quantity, Price: price, Image: image}
}

func Puid() string {
	id := uuid.NewV4()
	return id.String()
}

func getProduct() ([]byte, error) {
	var db, _ = sql.Open("sqlite3", "cache/users.sqlite3")
	defer db.Close()
	// var pid, product, quantity, price, image string
	rows, err := db.Query("select * from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return nil, err
	}
	//  fmt.Println(string(jsonData))
		return jsonData, nil
}


/*
// Helper function UnmarshalJSONTuple unmarshals JSON list (tuple) into a struct.

func UnmarshalJSONTuple(text []byte, obj interface{}) (err error) {
    var list []json.RawMessage
    err = json.Unmarshal(text, &list)
    if err != nil {
        return
    }

    objValue := reflect.ValueOf(obj).Elem()
    if len(list) > objValue.Type().NumField() {
        return fmt.Errorf("tuple has too many fields (%v) for %v",
            len(list), objValue.Type().Name())
    }

    for i, elemText := range list {
        err = json.Unmarshal(elemText, objValue.Field(i).Addr().Interface())
        if err != nil {
            return
        }
    }
    return
}

*/

