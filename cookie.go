// package Webcookie for managing sessions
package main

// import package
import "net/http"

// function SetSession creates a cookie for a user web session after successful user login 
func SetSession(u *User, w http.ResponseWriter) {
	value := map[string]string{
		"uuid": u.Uuid,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

// function GetUuid retrieves the uuid string from the cookie by decoding the "value" parameter of the cookie
func GetUuid(r *http.Request) (uuid string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			uuid = cookieValue["uuid"]
		}
	}
	return uuid
}

// function ClearSession clears the cookie created for the user web session
func ClearSession(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// function GetMsg is used decode the Alert message fron the Cookie "value" parameter and send it to the client application
func GetMsg(w http.ResponseWriter, r *http.Request, name string) (msg string) {
	if cookie, err := r.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			ClearSession(w, name)
		}
	}

	return msg
}

// function SetMsg is used to receive a message and create the alerts for the particular user web session. 
// This is done by encoding the message into the cookie value parameter. 
func SetMsg(w http.ResponseWriter, name string, msg string) {
	value := map[string]string{
		name: msg,
	}
	if encoded, err := cookieHandler.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

