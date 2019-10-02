package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func connect_db() {
	db, err = sql.Open("mysql", "root:10184902125410@/golang_db")

	if err != nil {
		log.Fatalln(err)
		log.Printf("Server:database wrong login and mysql password")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

type user struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
	Password  string
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Printf("Server: [net/http] method [%s] connection from [%v]", r.Method, r.RemoteAddr)

		next.ServeHTTP(w, r)
	}
}

func checkErrHandler(response http.ResponseWriter, Request *http.Request, err error) bool {
	if err != nil {

		fmt.Println(Request.Host + Request.URL.Path)

		//http.Redirect(response, Request, Request.Host+response.URL.Path, 301)
		return false
	}

	return true
}

func homeHandler(response http.ResponseWriter, Request *http.Request) {
	session := sessions.Start(response, Request)
	if len(session.GetString("username")) == 0 {
		http.Redirect(response, Request, "/login", 301)
	}

	var data = map[string]string{
		"username": session.GetString("username"),
		"message":  "Welcome to the Go !",
	}
	var t, err = template.ParseFiles("views/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(response, data)
	return

}

func QueryUserHandler(username string) user {
	var users = user{}
	err = db.QueryRow(`
		SELECT id, 
		username, 
		first_name, 
		last_name, 
		password 
		FROM users WHERE username=?
		`, username).
		Scan(
			&users.ID,
			&users.Username,
			&users.FirstName,
			&users.LastName,
			&users.Password,
		)
	return users
}

// Создание учетной записи
func signInPageHandler(response http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.Header().Set("Content-Type", "text/html; charset=utf-8")
		//http.ServerFile(response, Request, "views/signin.gohtml")
		return
	}

	username := r.FormValue("email")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	password := r.FormValue("password")

	users := QueryUserHandler(username)

	if (user{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErrHandler(response, r, err) {
			stmt, err := db.Prepare("INSERT INTO users SET username=?, password=?, first_name=?, last_name=?")
			if err == nil {
				_, err := stmt.Exec(&username, &hashedPassword, &first_name, &last_name)
				if err != nil {
					http.Error(response, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(response, r, "/login", http.StatusSeeOther)
				return
			}
		}
	} else {
		http.Redirect(response, r, "/register", 302)
	}

	/*
		err := json.NewDecoder(response.Body).Decode(&accLP)
		if err := nil {
			http.Redirect(response, Request, "/login", 301)
			log.Printf("Server:loginPage -> status redirect /login.[CODE:301]")
			return
		}

		err := bcrypt.CompareHashAndPassword([byte](accLP.databasePassword), [byte](password))
		if err := nil {
			http.Redirect(response, Request, "/login", 301)
			log.Printf("Server:loginPage -> status wrong password  redirect to /login.[CODE:301]")
			return
		}

		response.Write([]byte("Hello " + accLP.databaseUsername))
		log.Printf("Server greets user: ", accLP.databaseUsername)
		return

		response.WriteHeader(http.StatusBadRequest)
		log.Printf("Server:loginPage -> status bad reguest.")

		expectedPassword, ok := users[accLP.databaseUsername]

		if !ok || expectedPassword != accLP.databasePassword {
			response.WriteHeader(http.StatusUnauthorized)
			log.Printf("Server:loginPage -> status Unauthorized.")
			return
		}
	*/
}

func signUpHandler(response http.ResponseWriter, r *http.Request) {
	session := sessions.Start(response, r)
	if len(session.GetString("username")) != 0 && checkErrHandler(response, r, err) {
		http.Redirect(response, r, "/", 302)
	}
	if r.Method != "POST" {
		http.ServeFile(response, r, "views/login.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	users := QueryUserHandler(username)

	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

	if password_tes == nil {
		//login success
		session := sessions.Start(response, r)
		session.Set("username", users.Username)
		session.Set("name", users.FirstName)
		http.Redirect(response, r, "/", 302)
	} else {
		//login failed
		http.Redirect(response, r, "/login", 302)
	}

}

func logoutHandler(response http.ResponseWriter, Request *http.Request) {
	session := sessions.Start(response, Request)
	session.Clear()
	sessions.Destroy(response, Request)
	http.Redirect(response, Request, "/", 302)
}
