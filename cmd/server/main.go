// Copyright 2023 The Echo gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nxtaai/echo-grpc/api/echo/v1/echov1connect"
	"github.com/nxtaai/echo-grpc/echo"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(echov1connect.NewEchoAPIHandler(&echo.EchoAPIHandler{}))

	addr := "localhost:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = "localhost:" + port
	}
	log.Printf("Starting server on %s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 * 1024, // 8KB
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}
	keepAliveListener := tcpKeepAliveListener{ln.(*net.TCPListener)}

	go func() {
		err := srv.Serve(keepAliveListener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err)
	}
	log.Println("Server stopped")
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlivePeriod(30 * time.Second); err != nil {
		return nil, err
	}
	return tc, nil
}
