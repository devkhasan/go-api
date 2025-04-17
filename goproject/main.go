package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type Article struct {
	ID         int       `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Body       string    `db:"body" json:"body"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Author     string    `db:"author" json:"author"`
}

func main() {
	err := godotenv.Load() //loads from .env file
	if err != nil {
		log.Fatal("error loading .env file")
	}
	dbUrl := os.Getenv("DB_URL") // reads from the env
	log.Println("database URL:", dbUrl)

	//
	db, err = sqlx.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("couldnt open database", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database", err)

	}
	log.Println("succefully connected to database ")

	r := mux.NewRouter()
	r.HandleFunc("/", articles).Methods("GET")
	r.HandleFunc("/addArticles", add_article).Methods("POST")
	r.HandleFunc("/articles/{id}", delete).Methods("DELETE")
	r.HandleFunc("/changearrs/{id}", change).Methods("PUT")
	log.Println("Server started in :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
func articles(w http.ResponseWriter, r *http.Request) {

	var articles []Article
	err := db.Select(&articles, " SELECT * FROM articles")
	if err != nil {
		http.Error(w, "cannt  retrieve from database ", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(articles)

}
func add_article(w http.ResponseWriter, r *http.Request) {

	var newArticle Article
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		http.Error(w, " invalid request body", http.StatusBadRequest)
		return

	}
	_, err = db.Exec(`INSERT INTO articles(title,body,author) VALUES ($1,$2,$3)`,
		newArticle.Title, newArticle.Body, newArticle.Author)
	if err != nil {
		http.Error(w, "couldnt insert data ", http.StatusInternalServerError)

	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Article inserted successfully"})

}
func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec(`DELETE FROM articles WHERE ID=($1)`, id)
	if err != nil {
		http.Error(w, "couldnt delte from database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Article deleted successfully"})

}
func change(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updated Article
	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "bad bad rqueasr body", http.StatusBadRequest)
		return
	}
	_, err = db.Exec(`UPDATE articles SET title=$1, body=$2, author=$3 WHERE id=$4`,
		updated.Title, updated.Body, updated.Author, id)
	if err != nil {
		http.Error(w, "can cahnge the value", http.StatusInternalServerError)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Article updated successfully"})

}
