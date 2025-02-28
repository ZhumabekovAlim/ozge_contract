package main

import (
	"crypto/tls"
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
	defer db.Close()

	app := initializeApp(db, errorLog, infoLog)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           600,
	})

	// HTTP → HTTPS редирект
	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
		if err != nil {
			log.Fatal("HTTP redirect server error:", err)
		}
	}()

	// Настройки TLS
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// TCP-листенер с Keep-Alive
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		errorLog.Fatal(err)
	}
	tcpListener, ok := ln.(*net.TCPListener)
	if !ok {
		errorLog.Fatal("listener не является TCPListener")
	}
	listener := tcpKeepAliveListener{tcpListener}

	// HTTP сервер
	srv := &http.Server{
		Addr:              *addr,
		ErrorLog:          errorLog,
		Handler:           c.Handler(app.routes()),
		IdleTimeout:       2 * time.Minute,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		TLSConfig:         tlsConfig,
	}
	srv.SetKeepAlivesEnabled(true)
	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

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

// Редирект HTTP → HTTPS с минимальной задержкой
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
	}
}

// Keep-Alive TCP
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute)
	return tc, nil
}

// Хэндлер здоровья сервера
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
