stop:
	docker-compose down

clean:
	docker-compose down -v 

build:
	docker-compose build

up:
	docker-compose up

linter:
	golangci-lint run

restart: stop clean build up