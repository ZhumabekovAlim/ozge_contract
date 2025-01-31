package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"ozge/internal/config"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	port := os.Getenv("PORT")
	if port != "" {
		port = ":" + port
	} else {
		port = ":4000"
	}

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.URL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	app := initializeApp(db, errorLog, infoLog)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешает запросы с любых доменов (небезопасно, лучше указывать конкретные)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,          // Если используете куки или JWT
		AllowedHeaders:   []string{"*"}, // Разрешить все заголовки
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           600, // Кэширование preflight-запросов (секунды)
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      c.Handler(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Пути к SSL-сертификатам (замените на свои)
	certFile := "/etc/ssl/myserver/selfsigned.crt"
	keyFile := "/etc/ssl/myserver/selfsigned.key"

	infoLog.Printf("Starting HTTPS server on %s", *addr)

	err = srv.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		errorLog.Fatal(err)
	}

	select {}
}
