package main

import (
	"context"
	"flag"
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

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	env := flag.String("env", "", ".env file path")
	if *env == "" {
		godotenv.Load()
	} else {
		godotenv.Load(*env)
	}

	db, err := gc.InitDB()
	if err != nil {
		panic(err)
	}

	g := gc.NewGC(db)

	c := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(g.Router)

	srv := &http.Server{
		Handler:      c,
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
