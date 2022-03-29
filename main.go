package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conectDB() (connect *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	Name := "crud"

	connect, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return connect
}

var templates = template.Must(template.ParseGlob(("templates/*")))

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)

	log.Print("Server running")
	http.ListenAndServe(":8080", nil)
}

type User struct {
	Id    int
	Name  string
	Email string
}

func Home(w http.ResponseWriter, r *http.Request) {
	connectionEstablished := conectDB()
	register, err := connectionEstablished.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}
	user := User{}
	newUser := []User{}

	for register.Next() {
		var id int
		var name, email string
		err = register.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Email = email

		newUser = append(newUser, user)
	}
	fmt.Println(newUser)

	// fmt.Fprintf(w, "Holi")
	templates.ExecuteTemplate(w, "home", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		connectionEstablished := conectDB()
		insertRegister, err := connectionEstablished.Prepare("INSERT INTO users(name, email) VALUES (?, ?)")

		if err != nil {
			panic(err.Error())
		}
		insertRegister.Exec(name, email)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}
