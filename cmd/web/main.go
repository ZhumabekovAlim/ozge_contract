package main

import (
	"compress/gzip"
	"crypto/tls"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
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

	//c := cors.New(cors.Options{})

	// HTTP → HTTPS редирект
	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
		if err != nil {
			log.Fatal("HTTP redirect server error:", err)
		}
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Настройки TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // Поддержка TLS 1.2 и выше
		MaxVersion: tls.VersionTLS13, // Добавляем поддержку TLS 1.3
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		SessionTicketsDisabled:   false, // Включаем session tickets для быстрого соединения
		PreferServerCipherSuites: false,
	}

	// TCP-листенер с Keep-Alive
	ln, err := net.Listen("tcp4", *addr)
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
		Handler:           app.routes(),
		IdleTimeout:       2 * time.Minute,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		TLSConfig:         tlsConfig,
	}
	srv.SetKeepAlivesEnabled(true)

	// Пути к SSL-сертификатам Let's Encrypt
	certFile := "/etc/letsencrypt/live/infosite.kz/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/infosite.kz/privkey.pem"

	infoLog.Printf("Starting HTTPS server on %s", *addr)

	tlsConfig.Time = func() time.Time { return time.Now().Add(60 * time.Second) }
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

type gzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gz.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.Header().Set("Content-Encoding", "gzip")
	w.ResponseWriter.WriteHeader(statusCode)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, поддерживает ли клиент gzip
		if r.Header.Get("Accept-Encoding") != "" && r.Header.Get("Accept-Encoding") != "gzip" {
			next.ServeHTTP(w, r)
			return
		}

		// Создаем gzip.Writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Оборачиваем ResponseWriter
		wrappedWriter := &gzipResponseWriter{ResponseWriter: w, gz: gz}

		next.ServeHTTP(wrappedWriter, r)
	})
}
