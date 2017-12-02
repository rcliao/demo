package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type comment struct {
	Name    string
	Comment string
}

func main() {
	r := mux.NewRouter()
	defaultProtocol := "tcp"
	defaultPort := "3306"

	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		"root",
		"",
		defaultProtocol,
		"localhost",
		defaultPort,
		"demo",
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/", Index(db)).Methods("GET", "HEAD")
	r.HandleFunc("/hello", hello).Methods("GET", "HEAD")
	r.HandleFunc("/create", createComment(db)).Methods("POST")
	r.PathPrefix("/static").Handler(static())

	log.Println("Running web server at port 8000")
	http.ListenAndServe(":8000", r)
}

// Index renders the index page for submitting SQL queries to test
func Index(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./frontend/index.html")
		if err != nil {
			panic(err)
		}
		rows, err := db.Query("SELECT name, comment FROM comments")
		if err != nil {
			panic(err)
		}
		result := []comment{}
		for rows.Next() {
			var c comment
			if err := rows.Scan(&c.Name, &c.Comment); err != nil {
				log.Fatal(err)
			}
			result = append(result, c)
		}
		t.Execute(w, result)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello")
}

func createComment(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		comment := r.FormValue("comment")
		_, err := db.Query("INSERT INTO comments (name, comment) VALUES ('" + username + "', '" + comment + "')")
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", 302)
	})
}

func static() http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir("./frontend/static")))
}
