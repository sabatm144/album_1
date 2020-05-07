package route

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/julienschmidt/httprouter"
)

type customWriter struct {
	http.ResponseWriter
	statusCode int
}

func (cw *customWriter) WriteHeader(code int) {
	cw.statusCode = code
	cw.ResponseWriter.WriteHeader(code)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		cw := &customWriter{w, http.StatusOK}
		next.ServeHTTP(cw, r)
		t2 := time.Now()
		log.Printf("[%s] %d %q %v\n", r.Method, cw.statusCode, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// Prevent abnormal shutdown while panic
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Println(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func wrapHandler(next http.Handler) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", params)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return fn
}
