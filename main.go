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

	"github.com/mauryasaurav/timescale_database/repository"

	db "github.com/mauryasaurav/timescale_database/database"
	"github.com/mauryasaurav/timescale_database/handlers"
	"github.com/mauryasaurav/timescale_database/usecases"

	ginzap "github.com/gin-contrib/zap"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var dbSchema = `
create extension if not exists "pgcrypto";
create schema if not exists public;
create table if not exists public.users (
	id uuid default public.gen_random_uuid() not null,
	username character varying,
	time timestamp with time zone not null,
	email character varying,
	req_bytes int not null,
	pass character varying not null,
	primary key(id, time)
);

SELECT create_hypertable('users', 'time', migrate_data => true, if_not_exists => true);

CREATE MATERIALIZED VIEW if not exists users_hourly( time, req_bytes, username)
WITH (timescaledb.continuous) AS
  SELECT 
  		time_bucket('1h', time), 
  		sum(req_bytes),
		username
    FROM users
    GROUP BY time_bucket('1h', time), username WITH NO DATA;
`

func main() {

	// Create context  that listens for the interrupt signal from the os
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Loading TOML file
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] loading toml file: %+v\n", err)
		return
	}

	// Connecting To The Postgres Database
	conn, err := db.ConnectTimescaleDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] initializing postgres db: %+v\n", err)
		return
	}

	// Migrate db table
	conn.MustExec(dbSchema)

	r := gin.Default()

	logger, _ := zap.NewProduction()

	/*  Add a ginzap middleware, which:
	    - Logs all requests, like a combined access and error log.
	    - Logs to stdout.
		- RFC3339 with UTC time format.
	*/
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Declare user handler for postgres
	repo := repository.NewUserRepository(conn)
	userUseCase := usecases.NewUserUseCase(conn, repo)
	handlers.NewUserHandler(r, userUseCase)

	// Added default port address
	port := ":5000"
	if portEnv := config.Get("env.port").(string); len(portEnv) != 0 {
		port = portEnv
	}

	// Start the server
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Running server on goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish. the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
