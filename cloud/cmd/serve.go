package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func Serve(srv *http.Server) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()
	// 监听退出信号
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Received signal, gracefully shutting down...")
		cancel()
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	}
	wg.Wait()
}

type loggingMiddleware struct {
}

func (l loggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer func() {
			log.Printf(
				"%s %s %s %v",
				r.Method,
				r.URL.Path,
				r.Proto,
				time.Since(startTime),
			)
		}()
		next.ServeHTTP(w, r)
	})
}
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// log.Printf("Recovered from panic: %v", err)
				// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				ErrorMsg(w, "Internal Server Error:"+fmt.Sprintf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS 中间件
func corsMiddleware() mux.MiddlewareFunc {
	allowHeaders := "Content-Type, Accept, Authorization, Origin,Pwd"
	allowMethods := "GET, POST, PUT, DELETE, OPTIONS"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

			// 如果是预检请求（OPTIONS），直接返回 200 OK
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
