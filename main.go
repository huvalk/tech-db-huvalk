package main

import (
	"github.com/gorilla/mux"
	httpHuvalk "github.com/huvalk/tech-db-huvalk/api/http"
	"github.com/huvalk/tech-db-huvalk/api/repository"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     5432,
			Database: os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		MaxConnections: 100,
	})

	if err != nil {
		log.Println("Failed to open db: ", err.Error())
		return
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)
	handlers := httpHuvalk.NewHandler(repo)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.HandleFunc("/service/status", handlers.ServiceGetStatus).Methods("GET")
	router.HandleFunc("/service/clear", handlers.ServiceClear).Methods("POST")

	router.HandleFunc("/user/{nickname}/create", handlers.UserCreate).Methods("POST")
	router.HandleFunc("/user/{nickname}/profile", handlers.UserGet).Methods("GET")
	router.HandleFunc("/user/{nickname}/profile", handlers.UserChange).Methods("POST")

	router.HandleFunc("/forum/create", handlers.ForumCreate).Methods("POST")
	router.HandleFunc("/forum/{slug}/details", handlers.ForumGet).Methods("GET")
	router.HandleFunc("/forum/{slug}/threads", handlers.ForumGetListOfThreads).Methods("GET")
	router.HandleFunc("/forum/{slug}/users", handlers.ForumGetListOfUsers).Methods("GET")

	router.HandleFunc("/forum/{slug}/create", handlers.ForumCreateThread).Methods("POST")
	router.HandleFunc("/thread/{slug_or_id}/details", handlers.ThreadGetDetales).Methods("GET")
	router.HandleFunc("/thread/{slug_or_id}/details", handlers.ThreadChange).Methods("POST")
	router.HandleFunc("/thread/{slug_or_id}/vote", handlers.ThreadVote).Methods("POST")

	router.HandleFunc("/thread/{slug_or_id}/create", handlers.PostsCreate).Methods("POST")
	router.HandleFunc("/thread/{slug_or_id}/posts", handlers.ThreadGetListOfPost).Methods("GET")
	router.HandleFunc("/post/{id}/details", handlers.PostChangeDetails).Methods("POST")
	router.HandleFunc("/post/{id}/details", handlers.PostGetDetails).Methods("GET")

	http.Handle("/", router)
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		log.Println(err)
	}
}
