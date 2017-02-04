package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

func indexPage(w http.ResponseWriter, r *http.Request) {
	msg := getMsg(w, r, "message")
	var u = &User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		render(w, "signin", u)
	} else {
		u := &User{}
		render(w, "signin", u)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	pass := r.FormValue("password")
	u := &User{Username: name, Password: pass}
	redirect := "/"
	if name != "" && pass != "" {
		if b, uuid := userExists(u); b == true {
			setSession(&User{Uuid: uuid}, w)
			if name != "admin" {
				redirect = "/example"
			} else {
				redirect = "/admin"
			}
		} else {
			setMsg(w, "message", "please signup or enter a valid username and password!")
		}
	} else {
		setMsg(w, "message", "Username or Password field are empty!")
	}
	http.Redirect(w, r, redirect, 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w, "session")
	http.Redirect(w, r, "/", 302)
}

func examplePage(w http.ResponseWriter, r *http.Request) {
	uuid := getUuid(r)
	U := getUserFromUuid(uuid)
	if uuid != "" {
		pdata, err := getProduct()
		if err != nil {
			fmt.Println(err)
		}

		type Usr struct {
			Id      string
			Usrname string
			Pwd     string
			Fnm     string
			Lnm     string
			Eml     string
		}

		type Prdata struct {
			Puid     string `json:"puid"`
			Pname    string `json:"pname"`
			Quantity string `json:"quantity"`
			Price    string `json:"price"`
			Image    string `json:"image"`
		}

		type Viewdata struct {
			Usr Usr
			Prd []Prdata
		}

		var Vd Viewdata

		Vd.Usr.Id = U.Uuid
		Vd.Usr.Usrname = U.Username
		Vd.Usr.Pwd = U.Password
		Vd.Usr.Fnm = U.Fname
		Vd.Usr.Lnm = U.Lname
		Vd.Usr.Eml = U.Email

		// Pr := make([]Prdata, 0)
		var Pr []Prdata

		err = json.Unmarshal(pdata, &Pr)
		if err != nil {
			fmt.Println(err)
		}

		Vd.Prd = Pr

		render(w, "internal", Vd)

	} else {
		setMsg(w, "message", "Please login first!")
		http.Redirect(w, r, "/", 302)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"] = getMsg(w, r, "lname")
		u.Errors["fname"] = getMsg(w, r, "fname")
		u.Errors["email"] = getMsg(w, r, "email")
		u.Errors["username"] = getMsg(w, r, "username")
		u.Errors["password"] = getMsg(w, r, "password")
		render(w, "signup", u)
	case "POST":
		if n := checkUser(r.FormValue("userName")); n == true {
			setMsg(w, "username", "User already exists. Please enter a unique username!")
			http.Redirect(w, r, "/signup", 302)
			return
		}
		u := &User{
			Uuid:     Uuid(),
			Fname:    r.FormValue("fName"),
			Lname:    r.FormValue("lName"),
			Email:    r.FormValue("email"),
			Username: r.FormValue("userName"),
			Password: r.FormValue("password"),
		}
		result, err := govalidator.ValidateStruct(u)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Lname"); re == true {
				setMsg(w, "lname", "Please enter a valid Last Name")
			}
			if re := strings.Contains(e, "Email"); re == true {
				setMsg(w, "email", "Please enter a valid Email Address!")
			}
			if re := strings.Contains(e, "Fname"); re == true {
				setMsg(w, "fname", "Please enter a valid First Name")
			}
			if re := strings.Contains(e, "Username"); re == true {
				setMsg(w, "username", "Please enter a valid Username!")
			}
			if re := strings.Contains(e, "Password"); re == true {
				setMsg(w, "password", "Please enter a Password!")
			}

		}
		if r.FormValue("password") != r.FormValue("cpassword") {
			setMsg(w, "password", "The passwords you entered do not Match!")
			http.Redirect(w, r, "/signup", 302)
			return
		}

		if result == true {
			u.Password = encryptPass(u.Password)
			saveData(u)
			http.Redirect(w, r, "/", 302)
			return
		}
		http.Redirect(w, r, "/signup", 302)

	}
}

func admin(w http.ResponseWriter, r *http.Request) {
	uuid := getUuid(r)
	u := getUserFromUuid(uuid)
	if uuid != "" {
		render(w, "admin", u)
	} else {
		setMsg(w, "message", "Please login first!")
		http.Redirect(w, r, "/", 302)
	}
}

func addproduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uuid := getUuid(r)
		p := &Product{}
		p.Errors = make(map[string]string)
		p.Errors["pname"] = getMsg(w, r, "pname")
		p.Errors["quantity"] = getMsg(w, r, "quantity")
		p.Errors["price"] = getMsg(w, r, "price")
		p.Errors["image"] = getMsg(w, r, "image")
		if uuid != "" {
			render(w, "addproduct", p)
		} else {
			setMsg(w, "message", "Your Session has expired. Please login again!")
			http.Redirect(w, r, "/", 302)
		}

	case "POST":
		if n := checkProduct(r.FormValue("productName")); n == true {
			setMsg(w, "pname", "Product already exists!")
			http.Redirect(w, r, "/addproduct", 302)
			return
		}

		var pimage string

		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		file, header, err := r.FormFile("productimage")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			pimage = "NULL"
		} else {
			pimage = govalidator.ToString(header.Filename)
		}

		p := &Product{
			Puid:     Puid(),
			Pname:    r.FormValue("productName"),
			Quantity: r.FormValue("quantity"),
			Price:    r.FormValue("price"),
			Image:    pimage,
		}

		result, err := govalidator.ValidateStruct(p)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Pname"); re == true {
				setMsg(w, "pname", "Please enter a valid Product Name")
			}
			if re := strings.Contains(e, "Quantity"); re == true {
				setMsg(w, "quantity", "Please enter a valid Quantity!")
			}
			if re := strings.Contains(e, "Price"); re == true {
				setMsg(w, "price", "Please enter a valid Price")
			}
			if re := strings.Contains(e, "Image") || pimage == "NULL"; re == true {
				setMsg(w, "image", "Please enter a valid Image!")
			}
		}

		if result == true {
			f, err := os.Create("./files/" + p.Image)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if _, err := io.Copy(f, file); err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			saveProductData(p)
			http.Redirect(w, r, "/admin", 302)
			return
		}

		http.Redirect(w, r, "/addproduct", 302)

	}
}

func render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("view/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	govalidator.SetFieldsRequiredByDefault(true)
	http.Handle("/initializr/", http.StripPrefix("/initializr/", http.FileServer(http.Dir("initializr"))))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/example", examplePage)
	router.HandleFunc("/signup", signup).Methods("POST", "GET")
	router.HandleFunc("/admin", admin)
	router.HandleFunc("/addproduct", addproduct).Methods("POST", "GET")
	http.Handle("/", router)
	http.ListenAndServe(":8090", nil)
}
