package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"log"
	"net"
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
		port = ":443"
	}

	addr := flag.String("addr", port, "HTTPS network address")
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

	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		}))
		if err != nil {
			return
		}
	}()

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS13, // Используем минимально TLS 1.3
		PreferServerCipherSuites: true,             // Отдаем предпочтение серверным шифрам
		SessionTicketsDisabled:   false,            // Включаем поддержку сессионных тикетов
		CipherSuites: []uint16{ // Оптимизированные шифры
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
		},
		CurvePreferences: []tls.CurveID{ // Эффективные эллиптические кривые
			tls.X25519,
			tls.CurveP256,
		},
	}

	// Создаем TCP-листенер с заданным адресом
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		errorLog.Fatal(err)
	}
	// Убеждаемся, что listener является TCPListener-ом
	tcpListener, ok := ln.(*net.TCPListener)
	if !ok {
		errorLog.Fatal("listener не является TCPListener")
	}
	// Оборачиваем listener для установки TCP keep-alive
	listener := tcpKeepAliveListener{tcpListener}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      c.Handler(app.routes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}
	srv.SetKeepAlivesEnabled(true)

	// Пути к SSL-сертификатам Let's Encrypt
	certFile := "/etc/letsencrypt/live/infosite.kz/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/infosite.kz/privkey.pem"

	infoLog.Printf("Starting HTTPS server on %s", *addr)

	err = srv.ServeTLS(listener, certFile, keyFile)
	if err != nil {
		errorLog.Fatal(err)
	}

	select {}

}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	err = tc.SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	err = tc.SetKeepAlivePeriod(3 * time.Minute)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
