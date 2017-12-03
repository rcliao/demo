package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type comment struct {
	Name    string
	Comment template.HTML
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
		comment = strings.Replace(
			comment,
			"cat",
			fmt.Sprintf("<img src=\"%s\">", getRandomCatPicture()),
			-1,
		)
		_, err := db.Query(
			"INSERT INTO comments (name, comment) VALUES (?, ?)",
			username,
			comment,
		)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", 302)
	})
}

type giphyResp struct {
	Data giphyData `json:"data"`
}
type giphyData struct {
	ImageURL string `json:"image_url"`
}

func getRandomCatPicture() string {
	// https://api.giphy.com/v1/gifs/random?api_key=dc6zaTOxFJmzC&tag=cat
	resp, err := http.Get("https://api.giphy.com/v1/gifs/random?api_key=dc6zaTOxFJmzC&tag=cat")
	if err != nil {
		fmt.Println("Failed to get response from Github", err)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println("Failed to read body from Github get token response", err)
		return ""
	}
	var respBody giphyResp
	err = json.Unmarshal(b, &respBody)
	if err != nil {
		fmt.Println("Failed to parse json body from Github get token response", err)
		return ""
	}
	return respBody.Data.ImageURL

}

func static() http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir("./frontend/static")))
}
