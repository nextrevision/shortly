default: local

test:
	go test -v ./...

build:
	{ test -d dist || mkdir dist; } && go build -v -o dist/shortly ./...

package:
	tar czf shortly.tgz -C dist shortly

local: local-infra server

server:
	SHORTLY_DBUSER=shortly \
	SHORTLY_DBPASS=shortly \
	SHORTLY_DBHOST=localhost \
	SHORTLY_DBNAME=shortly \
	SHORTLY_DBSSLMODE=disable \
	SHORTLY_MEMCACHESERVERS="localhost:11211" \
	go run ./...

local-infra:
	docker-compose up -d

shorten-url:
	curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' \
	--data-urlencode "url=http://localhost:8000/echo/$(openssl rand -base64 12)" \
	http://localhost:8000