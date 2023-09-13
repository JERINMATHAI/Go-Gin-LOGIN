package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

type DetailsFromUser struct {
	password, username string
}
type DetailsIn struct {
	password, username string
}
type notValid struct {
	Not string
}

// var usersDetails = map[string]DetailsFromUser{} // for user password and username
var SessionDetails = map[string]string{} // for session id and user id
//var u DetailsFromUser

var enteredDetails DetailsFromUser
var BuiltInDetails DetailsIn

func login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		// checking for existing cookies
		if _, ok := SessionDetails[enteredDetails.username]; ok {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			// if there is not cookie open login page
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	} else if r.Method == "POST" {

		BuiltInDetails.username = "abcd"
		BuiltInDetails.password = "abcd"

		enteredDetails.username = r.FormValue("username")
		enteredDetails.password = r.FormValue("password")

		if enteredDetails.username == BuiltInDetails.username && enteredDetails.password == BuiltInDetails.password {
			// if new user, setting cookie

			cookieLogin, err := r.Cookie("sessionLogin")
			if err != nil {
				loginId := uuid.NewV4()
				cookieLogin = &http.Cookie{
					Name:     "sessionLogin",
					Value:    loginId.String(),
					HttpOnly: true,
				}
				http.SetCookie(w, cookieLogin)

				// get cookies
				//un := enteredDetails.username
				//ps := enteredDetails.password
				//u = DetailsFromUser{un, ps}

			}

			SessionDetails[enteredDetails.username] = cookieLogin.Value

			fmt.Println(cookieLogin)
			http.Redirect(w, r, "/home", http.StatusSeeOther)

		} else {
			t, _ := template.ParseFiles("login.html")
			p := notValid{Not: "Inavalid username or password"}
			t.Execute(w, p)
		}

		fmt.Println("User name", enteredDetails.username)
		fmt.Println("password:", enteredDetails.password)

	}

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
	//http.Redirect(w, r, "/back", http.StatusSeeOther)

}

func deleteCookie(w http.ResponseWriter, r *http.Request) {

	delete(SessionDetails, enteredDetails.username)
	fmt.Println(SessionDetails)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/home", home)
	http.HandleFunc("/back", deleteCookie)
	// set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
