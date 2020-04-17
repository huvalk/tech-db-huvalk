package main

import (
	httpHuvalk "github.com/huvalk/tech-db-huvalk/api/http"
	"github.com/huvalk/tech-db-huvalk/api/repository"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
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

	router := echo.New()

	router.GET("/api/service/status", handlers.ServiceGetStatus)
	router.POST("/api/service/clear", handlers.ServiceClear)

	router.POST("/api/user/:nickname/create", handlers.UserCreate)
	router.GET("/api/user/:nickname/profile", handlers.UserGet)
	router.POST("/api/user/:nickname/profile", handlers.UserChange)

	router.POST("/api/forum/create", handlers.ForumCreate)
	router.GET("/api/forum/:slug/details", handlers.ForumGet)
	router.GET("/api/forum/:slug/threads", handlers.ForumGetListOfThreads)
	router.GET("/api/forum/:slug/users", handlers.ForumGetListOfUsers)

	router.POST("/api/forum/:slug/create", handlers.ForumCreateThread)
	router.GET("/api/thread/:slug_or_id/details", handlers.ThreadGetDetales)
	router.POST("/api/thread/:slug_or_id/details", handlers.ThreadChange)
	router.POST("/api/thread/:slug_or_id/vote", handlers.ThreadVote)

	router.POST("/api/thread/:slug_or_id/create", handlers.PostsCreate)
	router.GET("/api/thread/:slug_or_id/posts", handlers.ThreadGetListOfPost)
	router.POST("/api/post/:id/details", handlers.PostChangeDetails)
	router.GET("/api/post/:id/details", handlers.PostGetDetails)

	http.Handle("/", router)
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		log.Println(err)
	}
}
