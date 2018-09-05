package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type middleware struct {
	log *log.Logger
}

func newMiddleware() middleware {
	return middleware{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (m middleware) accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.log.Printf("ACCESS: [%s] %s %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

func (m middleware) panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.log.Printf("ERROR: [%s] %s %s err: %v\n",
					r.Method, r.RemoteAddr, r.URL.Path, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var defaultMiddleware = newMiddleware()

func SetLogger(out io.Writer) {
	defaultMiddleware.log.SetOutput(out)
}

func AccessLog(next http.Handler) http.Handler {
	return defaultMiddleware.accessLog(next)
}

func Panic(next http.Handler) http.Handler {
	return defaultMiddleware.panic(next)
}
