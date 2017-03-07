// package for the main Go file
package main

// import packages
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"strconv"
	"time"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/stripe/stripe-go"
	// "github.com/stripe/stripe-go/currency"
	"github.com/stripe/stripe-go/charge"
)

// initialise variable to store cookie object
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// initialise the router variable
var router = mux.NewRouter()


// function that initializes the entry point of the web app and renders the landing page
func IndexPage(w http.ResponseWriter, r *http.Request) {
	msg := GetMsg(w, r, "message")
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

// Login function to authenticate user and render internal web page
func Login(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	pass := r.FormValue("password")
	u := &User{Username: name, Password: pass}
	redirect := "/"
	if name != "" && pass != "" {
		if b, uuid := UserExists(u); b == true {
			SetSession(&User{Uuid: uuid}, w)
			if name != "admin" {
				redirect = "/example"
			} else {
				redirect = "/admin"
			}
		} else {
			SetMsg(w, "message", "please signup or enter a valid username and password!")
		}
	} else {
		SetMsg(w, "message", "Username or Password field are empty!")
	}
	http.Redirect(w, r, redirect, 302)
}

// function Logout is to facilitate end of user's web session and clear the user session cookies
func Logout(w http.ResponseWriter, r *http.Request) {
	ClearSession(w, "session")
	http.Redirect(w, r, "/", 302)
}

// func Examplepage is to retrieve data from the database and send the data object to be rendered after User login
func ExamplePage(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	U := GetUserFromUuid(uuid)
	if uuid != "" {
		pdata, err := GetProduct()
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
		
		fmt.Println(Vd)

		render(w, "internal", Vd)

	} else {
		SetMsg(w, "message", "Please login first!")
		http.Redirect(w, r, "/", 302)
	}
}

// function Signup is the accept user registration details and save it to the database.
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"] = GetMsg(w, r, "lname")
		u.Errors["fname"] = GetMsg(w, r, "fname")
		u.Errors["email"] = GetMsg(w, r, "email")
		u.Errors["username"] = GetMsg(w, r, "username")
		u.Errors["password"] = GetMsg(w, r, "password")
		render(w, "Signup", u)
	case "POST":
		if n := CheckUser(r.FormValue("userName")); n == true {
			SetMsg(w, "username", "User already exists. Please enter a unique username!")
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
				SetMsg(w, "lname", "Please enter a valid Last Name")
			}
			if re := strings.Contains(e, "Email"); re == true {
				SetMsg(w, "email", "Please enter a valid Email Address!")
			}
			if re := strings.Contains(e, "Fname"); re == true {
				SetMsg(w, "fname", "Please enter a valid First Name")
			}
			if re := strings.Contains(e, "Username"); re == true {
				SetMsg(w, "username", "Please enter a valid Username!")
			}
			if re := strings.Contains(e, "Password"); re == true {
				SetMsg(w, "password", "Please enter a Password!")
			}

		}
		if r.FormValue("password") != r.FormValue("cpassword") {
			SetMsg(w, "password", "The passwords you entered do not Match!")
			http.Redirect(w, r, "/signup", 302)
			return
		}

		if result == true {
			u.Password = EncryptPass(u.Password)
			SaveData(u)
			http.Redirect(w, r, "/", 302)
			return
		}
		http.Redirect(w, r, "/signup", 302)

	}
}

// function Admin renders the main view after admin user login.
// order details are retrieved from the database and send as a data object to the client
func Admin(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	// u := GetUserFromUuid(uuid)
	if uuid != "" {
		
		orw, err := GetOrders()
		if err!=nil {
			fmt.Println(err)
		}
		
		var or ([]map[string]interface{})
		
		err = json.Unmarshal(orw, &or)
		if err !=nil {
			fmt.Println(err)
		}
		
		fmt.Println(or)
		
		render(w, "Admin", or)
	} else {
		SetMsg(w, "message", "Please login first!")
		http.Redirect(w, r, "/", 302)
	}
}

// function ManageProduct provides the functionality to view products and links for carrying out CRUD operations
func Manageproduct(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	// u := GetUserFromUuid(uuid)
	if uuid != "" {
		pdata, err := GetProduct()
		if err != nil {
			fmt.Println(err)
		}
		
		var pr ([]map[string]interface{})

		err = json.Unmarshal(pdata, &pr)
		if err != nil {
			fmt.Println(err)
		}
		
		fmt.Println(pr)
		
		render(w, "Manageproduct", pr)
	} else {
		SetMsg(w, "message", "Your session has expired. Please login again!")
		http.Redirect(w, r, "/", 302)
	}
}	

// function Addproduct renders the add product form, accepts admin users input and saves product data to the database
func Addproduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uuid := GetUuid(r)
		p := &Product{}
		p.Errors = make(map[string]string)
		p.Errors["pname"] = GetMsg(w, r, "pname")
		p.Errors["quantity"] = GetMsg(w, r, "quantity")
		p.Errors["price"] = GetMsg(w, r, "price")
		p.Errors["image"] = GetMsg(w, r, "image")
		
		fmt.Println("Getting to AddProducts...")
		
		if uuid != "" {
			render(w, "Addproduct", p)
		} else {
			SetMsg(w, "message", "Your Session has expired. Please login again!")
			http.Redirect(w, r, "/", 302)
		}

	case "POST":
	
		fmt.Println("Post Request sent Sucessfully...")
	
		if n := CheckProduct(r.FormValue("productName")); n == true {
			SetMsg(w, "pname", "Product already exists!")
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
		
		fmt.Println("Read file Sucessfully...")

		

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
				SetMsg(w, "pname", "Please enter a valid Product Name")
			}
			if re := strings.Contains(e, "Quantity"); re == true {
				SetMsg(w, "quantity", "Please enter a valid Quantity!")
			}
			if re := strings.Contains(e, "Price"); re == true {
				SetMsg(w, "price", "Please enter a valid Price")
			}
			if re := strings.Contains(e, "Image") || pimage == "NULL"; re == true {
				SetMsg(w, "image", "Please enter a valid Image!")
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
			
			fmt.Println("Saved File Image...")

			SaveProductData(p)
			
			fmt.Println("Saved Product Data...")
			
			http.Redirect(w, r, "/admin", 302)
			return
		}

		http.Redirect(w, r, "/addproduct", 302)

	}
}

//
func Showproduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uuid := GetUuid(r)
		fmt.Println(uuid)
		fmt.Println(r.FormValue("prodid"))
		pid := r.FormValue("prodid")
		fmt.Println(pid)
		a := GetProductFromPuid(pid)
		a.Errors = make(map[string]string)
		a.Errors["pname"] = GetMsg(w, r, "pname")
		a.Errors["quantity"] = GetMsg(w, r, "quantity")
		a.Errors["price"] = GetMsg(w, r, "price")
		a.Errors["image"] = GetMsg(w, r, "image")
		fmt.Println(a)
		if uuid != "" {
			render(w, "Showproduct", a)
		} else {
			SetMsg(w, "message", "Your Session has expired. Please login again!")
			http.Redirect(w, r, "/", 302)
		}

	case "POST":
	
		fmt.Println("Post Request sent Sucessfully...")
	
		if n := CheckProduct(r.FormValue("productName")); n == true {
			SetMsg(w, "pname", "Product already exists!")
			http.Redirect(w, r, "/showproduct", 302)
			return
		}

		var pimage string
		pexists := 1

		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		file, header, err := r.FormFile("productimage")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			pimage =  r.FormValue("imgname")
		} else {
			pimage = govalidator.ToString(header.Filename)
			pexists = 2
		}
		
		fmt.Println("Read file Sucessfully...")
		
		pu := r.FormValue("prodid")
		
		b := GetProductFromPuid(r.FormValue("prodid"))
		fmt.Println(b.Pname)
		
		b.Puid = pu

		a := &Product{
			Puid:     pu,
			Pname:    r.FormValue("productName"),
			Quantity: r.FormValue("quantity"),
			Price:    r.FormValue("price"),
			Image:    pimage,
		}
		
		fmt.Println(a)
		a.Errors = make(map[string]string)
		
		result, err := govalidator.ValidateStruct(a)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Pname"); re == true {
				// SetMsg(w, "pname", "Please enter a valid Product Name")
				a.Errors["pname"] = "Please enter a valid Product Name"
				a.Pname = b.Pname
			}
			if re := strings.Contains(e, "Quantity"); re == true {
				// SetMsg(w, "quantity", "Please enter a valid Quantity!")
				a.Errors["quantity"] = "Please enter a valid Quantity!"
				a.Quantity = b.Quantity
			}
			if re := strings.Contains(e, "Price"); re == true {
				// SetMsg(w, "price", "Please enter a valid Price")
				a.Errors["price"] = "Please enter a valid Price"
				a.Price = b.Price
			}
			if re := strings.Contains(e, "Image") || pimage == "NULL"; re == true {
				// SetMsg(w, "image", "Please enter a valid Image!")
				a.Errors["image"] = "Please enter a valid Image!"
			}
		}
		
		fmt.Println(result)
		
		if result == true {
			if pexists == 2 {
				fmt.Println(pexists)
				f, err := os.Create("./files/" + a.Image)
				if err != nil {
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer f.Close()
				if _, err := io.Copy(f, file); err != nil {
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			
				fmt.Println("Saved File Image...")
			}
			
			fmt.Println(a.Price)
			
			UpdateProductData(a)

			fmt.Println("Updated Product Data...")
			
			http.Redirect(w, r, "/manageproduct", 302)
			return
		}
		
		// http.Redirect(w, r, "/showproduct", 302)
		
		uuid := GetUuid(r)
		if uuid != "" {
			render(w, "Showproduct", a)
		} else {
			SetMsg(w, "message", "Your Session has expired. Please login again!")
			http.Redirect(w, r, "/", 302)
		}
		
	}		
}


func Deleteproduct(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	if uuid != "" {
		pid := r.FormValue("proid")
		fmt.Println(pid)
		DeleteProductData(pid)
		fmt.Println("Deleting Product Entries ...")
		http.Redirect(w, r, "/manageproduct", 302)
	}
}

// function Vieworder retrieves from the database the single order data object and renders it on the admin page
func Vieworder(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	// u := GetUserFromUuid(uuid)
	if uuid != "" {
		oid := r.FormValue("orderid")
		fmt.Println(oid)
		or, err := GetOrderFromOuid(oid)
		if err !=nil {
			fmt.Println(err)
		}

		fmt.Println(or)
		
		render(w, "Vieworder", or)
		
	} else {
		SetMsg(w, "message", "Your session has expired. Please login again!")
		http.Redirect(w, r, "/", 302)
	}
}

// function Placeorder renders the place order page that displays the order basket and provides the link to make payment
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	pubkey := "pk_test_75gzALA1yuRnl1y7qWdimlov"
	if uuid != "" {
		render(w, "Placeorder", pubkey)
	} else {
		SetMsg(w, "message", "Your session has expired. Please login again!")
		http.Redirect(w, r, "/", 302)
	}
}

// function Payment provides the functionality to call Stripe's API for payment, collect shipping details, post charge details, 
// generate the order and save order details to the database 
// order details contain the product details and shipping details
func Payment(w http.ResponseWriter, r *http.Request) {
	uuid := GetUuid(r)
	if uuid != "" {
		emaddress := GetUserEmailFromUuid(uuid)
		
		stripe.Key = "sk_test_rmp1OSLAWeSj1cAbJ7CvG3Rl"
		token := r.FormValue("stripeToken")
		
		var m = make(map[string]interface{})
		
		m["eladdress"] = r.FormValue("stripeEmail")
		m["shaddressname"] = r.FormValue("stripeShippingName")
		m["shaddressline"] = r.FormValue("stripeShippingAddressLine1")
		m["shaddresszip"] = r.FormValue("stripeShippingAddressZip")
		// m["shaddressstate"] = r.FormValue("stripeShippingAddressState")
		m["shaddresscity"] = r.FormValue("stripeShippingAddressCity")
		m["shaddresscountry"] = r.FormValue("stripeShippingAddressCountry")
		
		
		// marshal data to json
		j, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
		}
				
		
		s := r.FormValue("totalprice")
		totalamount, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				panic(err)
		}
		
		// fmt.Println(Orderdata)		
		// fmt.Println(eladdress)
		// fmt.Println(shaddressname)
		// fmt.Println(shaddressline)
		// fmt.Println(shaddresszip)
		// fmt.Println(shaddresscity)
		// fmt.Println(shaddresscountry)
		// fmt.Println(totalamount)
		// fmt.Println(m)
		
		params := &stripe.ChargeParams{
			Amount: totalamount,
			Currency: "inr",
			Desc: "Order Payment",
		}
		params.SetSource(token)

		charge, err := charge.New(params)
		if err != nil {
			fmt.Println(err)
		}
		
		// fmt.Println(charge)
		// fmt.Println(charge.ID)
		// fmt.Println(charge.Status)
		
		tmst := strconv.FormatInt(time.Now().Unix(), 10)
		
		fmt.Println(tmst)

		O := &Order{
			Ouid:     			Ouid(),
			Userid:         	uuid,
			Token:				r.FormValue("stripeToken"),
			OrderDetail: 		r.FormValue("Output"),
			ShipDetail:			string(j),
			TotalAmount:        r.FormValue("totalprice"), 
			ChargedAmount :     charge.Amount,
			Chargeid:			charge.ID,
			ChargeStatus:		charge.Status,
			Timestamp:    		tmst,
		}		
		
		fmt.Println(O)
		
		// fmt.Println(time.Unix(tmst, 0).Format("02.01.2006 15:04:05"))
		
		if charge.Status == "succeeded" {
			err := SaveOrderData(O) 
			if err != nil {
				fmt.Println(err)
			} 
			
			OrderEmail(r.FormValue("stripeEmail"))
			if emaddress != r.FormValue("stripeEmail") {
				OrderEmail(emaddress)
			} 
			
			// fmt.Println(O.Ouid)
			or, err := GetOrderFromOuid(O.Ouid)
			if err !=nil {
					fmt.Println(err)
			}
			
			render(w, "Payment", or)
			
		} else {
			render(w, "Payment", "Transaction Failed. Please Try Again!")
		}
		
	} else {
		SetMsg(w, "message", "Your session has expired. Please login again!")
		http.Redirect(w, r, "/", 302)
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
	router.HandleFunc("/", IndexPage)
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/logout", Logout).Methods("POST")
	router.HandleFunc("/example", ExamplePage)
	router.HandleFunc("/signup", Signup).Methods("POST", "GET")
	router.HandleFunc("/admin", Admin)
	router.HandleFunc("/manageproduct", Manageproduct)
	router.HandleFunc("/addproduct", Addproduct).Methods("POST", "GET")
	router.HandleFunc("/showproduct", Showproduct)
	router.HandleFunc("/deleteproduct", Deleteproduct)	
	router.HandleFunc("/vieworder", Vieworder)
	router.HandleFunc("/placeorder", PlaceOrder)
	router.HandleFunc("/payment", Payment)
	http.Handle("/", router)
	// http.ListenAndServe(":8090", nil)
	go http.ListenAndServe(":8090", http.RedirectHandler("https://127.0.0.1:8091",301))
	err := http.ListenAndServeTLS(":8091", "tls/cert.pem", "tls/key.pem", nil)
	log.Fatal(err)
}
