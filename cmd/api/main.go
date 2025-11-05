package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diagnosis/luxsuv-api-v2/internal/app"
	"github.com/diagnosis/luxsuv-api-v2/internal/routes"
	"github.com/diagnosis/luxsuv-api-v2/internal/store"
	"github.com/diagnosis/luxsuv-api-v2/migrations"
)

func main() {
	_ = os.Getenv("APP_ENN")
	dsn := os.Getenv("DATABASE_URL")
	pool, err := store.OpenPool(dsn)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer pool.Close()
	if err = store.MigrateFS(dsn, migrations.FS, "."); err != nil {
		pool.Close()
		log.Fatalf("migrate error: %v", err)
	}
	appl := app.NewApplication()
	r := routes.SetRouter(appl)

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	errch := make(chan error, 1)
	go func() {
		appURL := fmt.Sprintf("http://%s:%s", os.Getenv("APP_DOMAIN"), os.Getenv("APP_PORT"))
		log.Printf("Server is running on %f", appURL)
		errch <- serv.ListenAndServe()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("shutting down gracefully")
	case err := <-errch:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server stopped peacefully!")

}
