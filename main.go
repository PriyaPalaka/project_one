package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var db []Signup

type Signup struct {
	username string
	email    string
	password string
}

type CreateAccountFeedback struct {
	ErrorMsg   string
	SuccessMsg string
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	signupTmpl := template.Must(template.ParseFiles("templates/signup.html"))

	signupTmpl.Execute(w, nil)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	signinTmpl := template.Must(template.ParseFiles("templates/signin.html"))

	signinTmpl.Execute(w, nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	dashboardTmpl := template.Must(template.ParseFiles("templates/dashboard.html"))

	dashboardTmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello World!")
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method not supported")
		return
	}

	feedback := CreateAccountFeedback{
		ErrorMsg:   "Passwords doesn't match",
		SuccessMsg: "Account created!",
	}

	signupTmpl := template.Must(template.ParseFiles("templates/signup.html"))
	r.ParseForm()
	name := r.FormValue("name")
	email := r.FormValue("emailName")
	password := r.FormValue("passName")
	repassword := r.FormValue("Re-enterName")
	fmt.Println(name, email, password, repassword)

	if password != repassword {
		feedback.SuccessMsg = ""
		signupTmpl.Execute(w, feedback)
		return
	}

	for _, account := range db {
		if email == account.email {
			feedback.SuccessMsg = ""
			feedback.ErrorMsg = "Account already exists"
			signupTmpl.Execute(w, feedback)
			return
		}
	}

	newAccount := Signup{
		username: name,
		email:    email,
		password: password,
	}
	db = append(db, newAccount)

	feedback.ErrorMsg = ""
	signupTmpl.Execute(w, feedback)
}

func loginAccountHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintf(w, "Method Not Found")
		return
	}

	feedback := CreateAccountFeedback{
		ErrorMsg:   "Invalid Login Credentials",
		SuccessMsg: "Login Successful!",
	}
	dashboardTmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	signinTmpl := template.Must(template.ParseFiles("templates/signin.html"))
	r.ParseForm()
	email := r.FormValue("emailName")
	password := r.FormValue("passwordName")
	fmt.Println(email, password)

	emailFound := false
	for _, accountExists := range db {
		if email == accountExists.email {
			emailFound = true
			if password == accountExists.password {
				feedback.ErrorMsg = ""
				// signinTmpl.Execute(w, feedback)
				userName := map[string]interface{}{"Username": accountExists.username}
				dashboardTmpl.Execute(w, userName)
				// fmt.Fprintf(w, "Login sucessful")
				return
			} else {
				// fmt.Fprintf(w, "Invalid input")
				feedback.SuccessMsg = ""
				signinTmpl.Execute(w, feedback)
				return
			}

		}
	}

	if !emailFound {
		feedback.SuccessMsg = ""
		feedback.ErrorMsg = "Email Not Found"
		signinTmpl.Execute(w, feedback)

	}

}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/dashbord", dashboardHandler)
	http.HandleFunc("/createAccount", createAccountHandler)
	http.HandleFunc("/loginAccount", loginAccountHandler)
	fmt.Println("Server running on 40000")
	if err := http.ListenAndServe("0.0.0.0:40000", nil); err != nil {
		log.Fatal(err)

	}

}
