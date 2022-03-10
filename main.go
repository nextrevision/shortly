package main

import (
	"database/sql"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type Config struct {
	Debug           bool
	Port            int `default:"8000"`
	DbUser          string
	DbPass          string
	DbName          string
	DbHost          string
	DbPort          int    `default:"5432"`
	DbSSLMode       string `default:"verify-full"`
	MemcacheServers string
}

func initPgDataStore(config Config) (*PgDataStore, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Info("db connection successful")

	return NewPgDataStore(db), nil
}

func initMemcache(config Config) (*MemcacheDataStore, error) {
	servers := strings.Split(config.MemcacheServers, ",")

	client := memcache.New(servers...)

	if err := client.Ping(); err != nil {
		return nil, err
	}

	log.Info("memcache connection successful")

	return NewMemcachedDataStore(client), nil
}

func main() {
	// Parse environment config
	var config Config
	if err := envconfig.Process("shortly", &config); err != nil {
		log.Fatalf("could not process shortly config: %s", err)
	}

	// Enable debug logging
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// Create postgres datastore
	pgds, err := initPgDataStore(config)
	if err != nil {
		log.Fatalf("could not estasblish db connection: %s", err)
	}

	// Create service layer
	var service *UrlService
	if config.MemcacheServers == "" {
		service = NewUrlService(pgds)
	} else {
		mcds, err := initMemcache(config)
		if err != nil {
			log.Fatalf("could not establish memcached conection: %s", err)
		}
		service = NewUrlServiceWithCache(pgds, mcds)
	}

	// Create API layer
	api := NewApi(service)

	// Start up and serve the rest API
	srv := &http.Server{
		Handler:      api.Router,
		Addr:         fmt.Sprintf(":%d", config.Port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Info("listening on :8000")
	log.Fatal(srv.ListenAndServe())
}
