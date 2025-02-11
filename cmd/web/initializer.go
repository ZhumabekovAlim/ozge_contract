package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ozge/internal/handlers"
	"ozge/internal/repositories"
	"ozge/internal/services"
)

type application struct {
	errorLog          *log.Logger
	infoLog           *log.Logger
	tooHandler        *handlers.TOOHandler
	ipHandler         *handlers.IPHandler
	individualHandler *handlers.IndividualHandler
}

func initializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *application {

	tooRepo := &repositories.TOORepository{Db: db}
	tooService := &services.TOOService{Repo: tooRepo}
	tooHandler := &handlers.TOOHandler{Service: tooService}

	ipRepo := &repositories.IPRepository{Db: db}
	ipService := &services.IPService{Repo: ipRepo}
	ipHandler := &handlers.IPHandler{Service: ipService}

	individualRepo := &repositories.IndividualRepository{Db: db}
	individualService := &services.IndividualService{Repo: individualRepo}
	individualHandler := &handlers.IndividualHandler{Service: individualService}

	return &application{
		errorLog:          errorLog,
		infoLog:           infoLog,
		tooHandler:        tooHandler,
		ipHandler:         ipHandler,
		individualHandler: individualHandler,
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("%v", err)
		panic("failed to connect to database")
		return nil, err
	}
	db.SetMaxIdleConns(90)
	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
		panic("failed to ping the database")
		return nil, err
	}
	fmt.Println("successfully connected")

	return db, nil
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
