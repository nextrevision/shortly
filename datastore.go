package main

import (
	"database/sql"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	dataStoreErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "datastore_errors",
		Help: "Number of errors encountered by datastore",
	}, []string{"datastore"})
	cacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_hits",
		Help: "Number of cache hits by datastore",
	}, []string{"datastore"})
	cacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_misses",
		Help: "Number of cache misses by datastore",
	}, []string{"datastore"})
)

// The DataStore interface represents a simple storage format to save and resolve urls
type DataStore interface {
	SaveUrl(id string, url string) error
	ResolveUrl(id string) (string, error)
}

// Postgres DataStore struct and methods
type PgDataStore struct {
	db *sql.DB
}

func NewPgDataStore(db *sql.DB) *PgDataStore {
	return &PgDataStore{
		db: db,
	}
}

func (ds *PgDataStore) SaveUrl(id string, url string) error {
	stmt := `INSERT INTO urls (id, url) VALUES ($1, $2)`
	_, err := ds.db.Exec(stmt, id, url)
	if err != nil {
		dataStoreErrors.WithLabelValues("postgres").Inc()
	}
	return err
}

func (ds *PgDataStore) ResolveUrl(id string) (string, error) {
	rows, err := ds.db.Query(`SELECT url FROM urls WHERE id = $1`, id)
	if err != nil {
		dataStoreErrors.WithLabelValues("postgres").Inc()
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {
			dataStoreErrors.WithLabelValues("postgres").Inc()
		}

		// Get first match
		return url, err
	}

	return "", nil
}

// Memcache DataStore struct and methods
type MemcacheDataStore struct {
	client *memcache.Client
	expiry float64
}

func NewMemcachedDataStore(client *memcache.Client) *MemcacheDataStore {
	return &MemcacheDataStore{
		client: client,
		expiry: (12 * time.Hour).Seconds(),
	}
}

func (ds *MemcacheDataStore) SaveUrl(id string, url string) error {
	if err := ds.client.Add(&memcache.Item{
		Key:        id,
		Value:      []byte(url),
		Expiration: int32(ds.expiry),
	}); err != nil {
		dataStoreErrors.WithLabelValues("memcache").Inc()
		return err
	}

	return nil
}

func (ds *MemcacheDataStore) ResolveUrl(id string) (string, error) {
	item, err := ds.client.Get(id)
	if err != nil {
		dataStoreErrors.WithLabelValues("memcache").Inc()
		return "", err
	}

	url := string(item.Value)
	if url == "" {
		cacheMisses.WithLabelValues("memcache").Inc()
	} else {
		cacheHits.WithLabelValues("memcache").Inc()
	}

	return url, nil
}

type MemoryDataStore struct {
	urls map[string]string
}

// Memory DataStore struct and methods; should only be used for testing
func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		urls: map[string]string{},
	}
}

func (ds *MemoryDataStore) SaveUrl(id string, url string) error {
	ds.urls[id] = url
	return nil
}

func (ds *MemoryDataStore) ResolveUrl(id string) (string, error) {
	return ds.urls[id], nil
}
