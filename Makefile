default: local

test:
	go test -v ./...

build:
	docker build -t nextrevision/shortly:latest .

publish:
	docker push nextrevision/shortly:latest

deploy: build publish

server: local-infra
	SHORTLY_DBUSER=shortly \
	SHORTLY_DBPASS=shortly \
	SHORTLY_DBHOST=localhost \
	SHORTLY_DBNAME=shortly \
	SHORTLY_DBSSLMODE=disable \
	SHORTLY_MEMCACHESERVERS="localhost:11211" \
	go run ./...

local: local-infra local-docker

local-docker: local-infra
	docker-compose up -d --build shortly

local-infra:
	docker-compose up -d memcached postgres liquibase-postgres

shorten-url:
	curl -X POST -H 'Content-Type: application/x-www-form-urlencoded' \
	--data-urlencode "url=http://localhost:8000/echo/test-`openssl rand -base64 12`" \
	http://localhost:8000

tf-plan:
	cd terraform && terraform init && terraform plan

tf-apply:
	cd terraform && terraform apply -auto-approve
