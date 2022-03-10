package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

const (
	urlChars              = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idRandLength          = 10
	idDeterministicLength = 6
)

// The UrlService acts as a proxy between datastore implementations to provide ease in adding layers and
//switching implementations
type UrlService struct {
	datastore    DataStore
	cache        DataStore
	cacheEnabled bool
}

func NewUrlService(ds DataStore) *UrlService {
	return &UrlService{
		datastore: ds,
	}
}

func NewUrlServiceWithCache(ds DataStore, cache DataStore) *UrlService {
	service := NewUrlService(ds)
	service.cacheEnabled = true
	service.cache = cache
	return service
}

func (s *UrlService) SaveUrl(url string) (string, error) {
	id, err := generateUrlId(url)
	if err != nil {
		return "", err
	}

	err = s.datastore.SaveUrl(id, url)
	if err != nil {
		return "", err
	}

	if s.cacheEnabled {
		err = s.cache.SaveUrl(id, url)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (s *UrlService) ResolveUrl(id string) (string, error) {
	var url string
	var err error

	if s.cacheEnabled {
		url, _ = s.cache.ResolveUrl(id)
	}

	if url == "" {
		url, err = s.datastore.ResolveUrl(id)
		if err != nil {
			return "", err
		}
	}

	return url, nil
}

// Basic algorithm to create a semi-random, semi-deterministic ID based on a url
func generateUrlId(url string) (string, error) {
	randBytes := make([]byte, idRandLength)
	rand.Seed(time.Now().UnixNano())
	for i := range randBytes {
		randBytes[i] = urlChars[rand.Intn(len(urlChars))]
	}

	urlEncoded := base64.StdEncoding.EncodeToString([]byte(url))
	urlIdSuffix := string([]byte(urlEncoded)[:idDeterministicLength])
	return fmt.Sprintf("%s%s", randBytes, urlIdSuffix), nil
}
