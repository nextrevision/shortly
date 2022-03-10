package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceSaveAndResolveUrl(t *testing.T) {
	url := "http://foo"
	ds := NewMemoryDataStore()
	service := NewUrlService(ds)

	id, err := service.SaveUrl(url)

	assert.Nil(t, err, "should not throw an error")
	assert.NotEmpty(t, id, "id should not be empty")

	resolved, err := ds.ResolveUrl(id)

	assert.Nil(t, err, "should not throw an error")
	assert.Equal(t, url, resolved, "urls should match")
}

func TestGenerateUrlId(t *testing.T) {
	id, err := generateUrlId("foo")
	assert.Nil(t, err, "should not throw an error")
	assert.NotEmpty(t, id, "id should not be empty")
}
