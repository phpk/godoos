/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"context"
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
