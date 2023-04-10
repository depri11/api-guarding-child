package main

import (
	"context"
	"fmt"
	gc "gc/src"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting...")

	db, err := gc.InitDB()
	if err != nil {
		panic(err)
	}

	g := gc.NewGC(db)

	srv := &http.Server{
		Handler:      g.Router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Println("API started in localhost", ":8080")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		done <- true
	}()
	<-done

	log.Println("server stoped")
	srv.Shutdown(context.Background())
}
