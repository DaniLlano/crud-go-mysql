package main

import (
	"database/sql"
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
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

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
	registers, err := connectionEstablished.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}
	user := User{}
	newUser := []User{}

	for registers.Next() {
		var id int
		var name, email string
		err = registers.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Email = email

		newUser = append(newUser, user)
	}

	templates.ExecuteTemplate(w, "home", newUser)
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

func Delete(w http.ResponseWriter, r *http.Request) {
	idUser := r.URL.Query().Get("id")

	connectionEstablished := conectDB()
	deleteRegister, err := connectionEstablished.Prepare("DELETE FROM users WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	deleteRegister.Exec(idUser)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idUser := r.URL.Query().Get("id")

	connectionEstablished := conectDB()
	register, err := connectionEstablished.Query("SELECT * FROM users WHERE id=?", idUser)

	user := User{}

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
	}

	templates.ExecuteTemplate(w, "update", user)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue(("id"))
		name := r.FormValue("name")
		email := r.FormValue("email")

		connectionEstablished := conectDB()
		updateRegister, err := connectionEstablished.Prepare("UPDATE users SET name=?, email=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		updateRegister.Exec(name, email, id)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}
