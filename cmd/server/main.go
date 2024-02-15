package main

import (
	"database/sql"
	"os"

	"github.com/codingconcepts/env"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/golden-vcr/auth"
	"github.com/golden-vcr/remix/gen/queries"
	"github.com/golden-vcr/remix/internal/admin"
	"github.com/golden-vcr/remix/internal/state"
	"github.com/golden-vcr/server-common/db"
	"github.com/golden-vcr/server-common/entry"
)

type Config struct {
	BindAddr   string `env:"BIND_ADDR"`
	ListenPort uint16 `env:"LISTEN_PORT" default:"5010"`

	AuthURL string `env:"AUTH_URL" default:"http://localhost:5002"`

	DatabaseHost     string `env:"PGHOST" required:"true"`
	DatabasePort     int    `env:"PGPORT" required:"true"`
	DatabaseName     string `env:"PGDATABASE" required:"true"`
	DatabaseUser     string `env:"PGUSER" required:"true"`
	DatabasePassword string `env:"PGPASSWORD" required:"true"`
	DatabaseSslMode  string `env:"PGSSLMODE"`
}

func main() {
	app, ctx := entry.NewApplication("remix")
	defer app.Stop()

	// Parse config from environment variables
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		app.Fail("Failed to load .env file", err)
	}
	config := Config{}
	if err := env.Set(&config); err != nil {
		app.Fail("Failed to load config", err)
	}

	// Initialize an auth client so we can require broadcaster-level access in order to
	// call admin-only API endpoints
	authClient, err := auth.NewClient(ctx, config.AuthURL)
	if err != nil {
		app.Fail("Failed to initialize auth client", err)
	}

	// Configure our database connection and initialize a Queries struct, so we can use
	// the 'remix' schema to store and retrieve persistent data related to clips,
	// playback, etc.
	connectionString := db.FormatConnectionString(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseSslMode,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		app.Fail("Failed to open sql.DB", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		app.Fail("Failed to connect to database", err)
	}
	q := queries.New(db)

	// Start setting up our HTTP handlers, using gorilla/mux for routing
	r := mux.NewRouter()

	// We can call the broadcaster-only admin API to register new clips etc.
	{
		adminServer := admin.NewServer(q)
		adminServer.RegisterRoutes(authClient, r.PathPrefix("/admin").Subrouter())
	}

	// Any client (no auth required) can call the state API to get read-only information
	// about available clips etc.
	{
		stateServer := state.NewServer(q)
		stateServer.RegisterRoutes(r)
	}

	// Handle incoming HTTP connections until our top-level context is canceled, at
	// which point shut down cleanly
	entry.RunServer(ctx, app.Log(), r, config.BindAddr, config.ListenPort)
}
